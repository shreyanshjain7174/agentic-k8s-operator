"""
Tests for MCP client SSRF protection.

Validates URL scheme checking, host blocking, and IP range filtering.
"""

import pytest
from unittest.mock import patch

from agents.tools.mcp_client import (
    MCPClient,
    SSRFConfig,
    SSRFProtectionError,
    URLValidationError,
    validate_url,
    validate_url_scheme,
    validate_url_host,
    is_url_allowed,
    is_private_ip,
    resolve_hostname,
)


class TestSchemeValidation:
    """Tests for URL scheme validation"""
    
    def test_valid_https_scheme(self):
        """HTTPS scheme should be allowed by default"""
        scheme = validate_url_scheme("https://example.com/api")
        assert scheme == "https"
    
    def test_reject_http_scheme(self):
        """HTTP scheme should be rejected by default"""
        with pytest.raises(URLValidationError) as exc:
            validate_url_scheme("http://example.com/api")
        assert "not allowed" in str(exc.value)
    
    def test_reject_ftp_scheme(self):
        """FTP scheme should be rejected"""
        with pytest.raises(URLValidationError):
            validate_url_scheme("ftp://example.com/file")
    
    def test_reject_file_scheme(self):
        """File scheme should be rejected"""
        with pytest.raises(URLValidationError):
            validate_url_scheme("file:///etc/passwd")
    
    def test_custom_allowed_schemes(self):
        """Should allow custom schemes when specified"""
        scheme = validate_url_scheme(
            "http://example.com/api",
            allowed_schemes={"http", "https"}
        )
        assert scheme == "http"
    
    def test_missing_scheme(self):
        """Should reject URLs without scheme"""
        with pytest.raises(URLValidationError) as exc:
            validate_url_scheme("example.com/api")
        assert "missing scheme" in str(exc.value)


class TestPrivateIPDetection:
    """Tests for private IP detection"""
    
    def test_loopback_ip(self):
        """127.0.0.1 should be detected as private"""
        assert is_private_ip("127.0.0.1") is True
    
    def test_10_range(self):
        """10.x.x.x should be detected as private"""
        assert is_private_ip("10.0.0.1") is True
        assert is_private_ip("10.255.255.255") is True
    
    def test_172_range(self):
        """172.16-31.x.x should be detected as private"""
        assert is_private_ip("172.16.0.1") is True
        assert is_private_ip("172.31.255.255") is True
    
    def test_192_range(self):
        """192.168.x.x should be detected as private"""
        assert is_private_ip("192.168.0.1") is True
        assert is_private_ip("192.168.255.255") is True
    
    def test_public_ip(self):
        """Public IPs should not be private"""
        assert is_private_ip("8.8.8.8") is False
        assert is_private_ip("1.1.1.1") is False
    
    def test_ipv6_loopback(self):
        """IPv6 loopback should be private"""
        assert is_private_ip("::1") is True
    
    def test_invalid_ip(self):
        """Invalid IPs should return False"""
        assert is_private_ip("not-an-ip") is False


class TestHostValidation:
    """Tests for host validation"""
    
    def test_localhost_blocked(self):
        """localhost should be blocked"""
        with pytest.raises(SSRFProtectionError) as exc:
            validate_url_host("https://localhost/api")
        assert "Localhost" in str(exc.value)
    
    def test_127_0_0_1_blocked(self):
        """127.0.0.1 should be blocked"""
        with pytest.raises(SSRFProtectionError) as exc:
            validate_url_host("https://127.0.0.1/api")
        assert "blocked" in str(exc.value).lower()
    
    def test_localhost_localdomain_blocked(self):
        """localhost.localdomain should be blocked"""
        with pytest.raises(SSRFProtectionError):
            validate_url_host("https://localhost.localdomain/api")
    
    def test_local_domain_blocked(self):
        """.local domains should be blocked"""
        with pytest.raises(SSRFProtectionError):
            validate_url_host("https://myservice.local/api")
    
    def test_allow_localhost_flag(self):
        """Should allow localhost when flag is set"""
        hostname = validate_url_host(
            "https://localhost/api",
            allow_localhost=True
        )
        assert hostname == "localhost"
    
    @patch("agents.tools.mcp_client.resolve_hostname")
    def test_private_ip_blocked(self, mock_resolve):
        """Private IPs should be blocked after DNS resolution"""
        mock_resolve.return_value = "10.0.0.50"
        
        with pytest.raises(SSRFProtectionError) as exc:
            validate_url_host("https://internal-service.example.com/api")
        assert "blocked range" in str(exc.value)
    
    @patch("agents.tools.mcp_client.resolve_hostname")
    def test_192_168_blocked(self, mock_resolve):
        """192.168.x.x should be blocked"""
        mock_resolve.return_value = "192.168.1.100"
        
        with pytest.raises(SSRFProtectionError):
            validate_url_host("https://internal.example.com/api")
    
    @patch("agents.tools.mcp_client.resolve_hostname")
    def test_172_range_blocked(self, mock_resolve):
        """172.16-31.x.x should be blocked"""
        mock_resolve.return_value = "172.16.5.5"
        
        with pytest.raises(SSRFProtectionError):
            validate_url_host("https://internal.example.com/api")
    
    def test_whitelist_validation(self):
        """Should validate against host whitelist"""
        with pytest.raises(SSRFProtectionError) as exc:
            validate_url_host(
                "https://malicious.com/api",
                allowed_hosts={"api.trusted.com", "api.safe.com"}
            )
        assert "not in whitelist" in str(exc.value)
    
    def test_whitelist_allows_valid_host(self):
        """Whitelist should allow valid hosts"""
        hostname = validate_url_host(
            "https://api.trusted.com/endpoint",
            allowed_hosts={"api.trusted.com", "api.safe.com"}
        )
        assert hostname == "api.trusted.com"
    
    def test_no_hostname(self):
        """Should reject URLs without hostname"""
        with pytest.raises(SSRFProtectionError):
            validate_url_host("https:///api")


class TestURLValidation:
    """Tests for full URL validation"""
    
    def test_valid_public_url(self):
        """Valid public HTTPS URLs should pass"""
        url = validate_url("https://api.example.com/endpoint")
        assert url == "https://api.example.com/endpoint"
    
    def test_reject_http(self):
        """Should reject HTTP URLs"""
        with pytest.raises(URLValidationError):
            validate_url("http://api.example.com/endpoint")
    
    def test_reject_localhost(self):
        """Should reject localhost"""
        with pytest.raises(SSRFProtectionError):
            validate_url("https://localhost/endpoint")
    
    def test_custom_config(self):
        """Should respect custom config"""
        config = SSRFConfig(
            allowed_schemes={"https", "http"},
            allow_localhost=True,
        )
        url = validate_url("http://localhost/test", config=config)
        assert url == "http://localhost/test"


class TestIsURLAllowed:
    """Tests for is_url_allowed helper"""
    
    def test_allowed_public_url(self):
        """Public HTTPS URL should be allowed"""
        assert is_url_allowed("https://api.example.com/endpoint") is True
    
    def test_blocked_http(self):
        """HTTP URL should be blocked by default"""
        assert is_url_allowed("http://api.example.com/endpoint") is False
    
    def test_blocked_localhost(self):
        """Localhost should be blocked"""
        assert is_url_allowed("https://localhost/endpoint") is False
    
    def test_whitelist_enforcement(self):
        """Should enforce whitelist"""
        assert is_url_allowed(
            "https://api.trusted.com/endpoint",
            allowed_hosts={"api.trusted.com"}
        ) is True
        
        assert is_url_allowed(
            "https://api.untrusted.com/endpoint",
            allowed_hosts={"api.trusted.com"}
        ) is False


class TestSSRFConfig:
    """Tests for SSRFConfig dataclass"""
    
    def test_default_config(self):
        """Default config should be secure"""
        config = SSRFConfig()
        assert config.allowed_schemes == {"https"}
        assert config.allow_localhost is False
        assert config.allow_private_ips is False
        assert len(config.blocked_ip_ranges) > 0
    
    def test_custom_schemes(self):
        """Should allow custom schemes"""
        config = SSRFConfig(allowed_schemes={"https", "wss"})
        assert config.allowed_schemes == {"https", "wss"}
    
    def test_custom_blocked_ranges(self):
        """Should allow custom blocked IP ranges"""
        config = SSRFConfig(blocked_ip_ranges=["10.0.0.0/8"])
        assert config.blocked_ip_ranges == ["10.0.0.0/8"]


class TestMCPClient:
    """Tests for MCPClient class"""
    
    def test_client_init_default_config(self):
        """Client should initialize with secure defaults"""
        client = MCPClient()
        assert client.config.allowed_schemes == {"https"}
        assert client.config.allow_localhost is False
    
    def test_client_init_custom_config(self):
        """Client should accept custom config"""
        config = SSRFConfig(allow_localhost=True)
        client = MCPClient(config=config)
        assert client.config.allow_localhost is True
    
    def test_validate_endpoint_valid(self):
        """Should validate valid endpoints"""
        client = MCPClient()
        url = client.validate_endpoint("https://api.example.com/mcp")
        assert url == "https://api.example.com/mcp"
    
    def test_validate_endpoint_blocked(self):
        """Should reject blocked endpoints"""
        client = MCPClient()
        with pytest.raises(SSRFProtectionError):
            client.validate_endpoint("https://localhost/mcp")
    
    @pytest.mark.asyncio
    async def test_context_manager(self):
        """Should work as async context manager"""
        async with MCPClient() as client:
            assert client._session is not None
    
    @pytest.mark.asyncio
    async def test_close_session(self):
        """Should close session properly"""
        client = MCPClient()
        await client._ensure_session()
        assert client._session is not None
        await client.close()
        assert client._session is None


class TestEdgeCases:
    """Tests for edge cases and security scenarios"""
    
    def test_ipv4_in_hostname(self):
        """Should handle IP addresses in hostname"""
        with pytest.raises(SSRFProtectionError):
            validate_url_host("https://10.0.0.1/api")
    
    def test_ipv6_localhost(self):
        """Should block IPv6 localhost"""
        with pytest.raises(SSRFProtectionError):
            validate_url_host("https://[::1]/api")
    
    def test_decimal_ip_encoding(self):
        """Should handle various URL encodings"""
        with pytest.raises((URLValidationError, SSRFProtectionError)):
            validate_url("https://2130706433/api")
    
    def test_hex_ip_encoding(self):
        """Should handle hex IP encodings"""
        with pytest.raises((URLValidationError, SSRFProtectionError)):
            validate_url("https://0x7f000001/api")
    
    def test_url_with_port(self):
        """Should handle URLs with ports"""
        url = validate_url("https://api.example.com:8443/endpoint")
        assert url == "https://api.example.com:8443/endpoint"
    
    def test_url_with_credentials(self):
        """Should handle URLs with credentials"""
        url = validate_url("https://user:pass@api.example.com/endpoint")
        assert "api.example.com" in url
    
    @patch("agents.tools.mcp_client.resolve_hostname")
    def test_dns_rebinding_protection(self, mock_resolve):
        """Should protect against DNS rebinding"""
        mock_resolve.return_value = "127.0.0.1"
        
        with pytest.raises(SSRFProtectionError):
            validate_url_host("https://attacker-controlled.com/api")
    
    @patch("agents.tools.mcp_client.resolve_hostname")
    def test_cloud_metadata_blocked(self, mock_resolve):
        """Should block cloud metadata endpoints"""
        mock_resolve.return_value = "169.254.169.254"
        
        with pytest.raises(SSRFProtectionError):
            validate_url_host("https://metadata.google.internal/computeMetadata/v1/")
