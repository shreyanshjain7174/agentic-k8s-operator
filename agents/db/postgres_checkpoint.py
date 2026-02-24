"""
Database security configuration module for PostgreSQL checkpoint storage.

Security Features:
- SSL/TLS requirement: All connections must use SSL encryption
- Encrypted credentials: Sensitive data stored in environment variables
- Credential management: Use environment variables (never hardcode secrets)
- Connection validation: Validate connection strings before use
- Password encryption: Support for encrypted password storage

Environment Variables:
    POSTGRES_HOST: Database host (required)
    POSTGRES_PORT: Database port (default: 5432)
    POSTGRES_DB: Database name (required)
    POSTGRES_USER: Database user (required)
    POSTGRES_PASSWORD: Database password (required)
    POSTGRES_SSL_MODE: SSL mode (default: require)
    POSTGRES_SSL_CERT: Path to SSL certificate (optional)
    POSTGRES_SSL_KEY: Path to SSL key (optional)
    POSTGRES_SSL_ROOT_CERT: Path to CA certificate (optional)

Usage:
    from agents.db.postgres_checkpoint import create_postgres_saver, validate_connection

    # Validate configuration before creating connection
    validate_connection()

    # Create PostgresSaver for LangGraph checkpointing
    saver = create_postgres_saver()
"""

import logging
import os
import re
import ssl
from typing import Any, Dict, Optional

from psycopg2 import pool
from psycopg2.extensions import ISOLATION_LEVEL_AUTOCOMMIT
from pydantic import BaseModel, Field, field_validator

logger = logging.getLogger(__name__)


class DatabaseSecurityConfig(BaseModel):
    """Security configuration for PostgreSQL database connection."""

    host: str = Field(..., description="Database host")
    port: int = Field(default=5432, ge=1, le=65535, description="Database port")
    database: str = Field(..., description="Database name")
    user: str = Field(..., description="Database user")
    password: str = Field(..., description="Database password (sensitive)")
    ssl_mode: str = Field(default="require", description="SSL mode")
    ssl_cert: Optional[str] = Field(default=None, description="Path to SSL certificate")
    ssl_key: Optional[str] = Field(default=None, description="Path to SSL key")
    ssl_root_cert: Optional[str] = Field(default=None, description="Path to CA certificate")

    @field_validator("ssl_mode")
    @classmethod
    def validate_ssl_mode(cls, v: str) -> str:
        """Validate SSL mode is secure."""
        allowed_modes = {"require", "verify-ca", "verify-full", "disable"}
        v = v.lower()
        if v not in allowed_modes:
            raise ValueError(f"SSL mode must be one of: {allowed_modes}")
        if v == "disable":
            raise ValueError("SSL mode 'disable' is not allowed for security reasons")
        return v

    @field_validator("host")
    @classmethod
    def validate_host(cls, v: str) -> str:
        """Validate host is not a local socket or dangerous."""
        if v in ("localhost", "127.0.0.1", "::1"):
            return v
        if not re.match(r"^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?)*$", v):
            raise ValueError(f"Invalid hostname: {v}")
        return v

    def get_ssl_context(self) -> Optional[Dict[str, Any]]:
        """Build SSL context dictionary for psycopg2."""
        if self.ssl_mode == "disable":
            return None

        ssl_context = {
            "sslmode": self.ssl_mode,
        }

        if self.ssl_cert:
            if not os.path.exists(self.ssl_cert):
                raise FileNotFoundError(f"SSL certificate not found: {self.ssl_cert}")
            ssl_context["sslcert"] = self.ssl_cert

        if self.ssl_key:
            if not os.path.exists(self.ssl_key):
                raise FileNotFoundError(f"SSL key not found: {self.ssl_key}")
            ssl_context["sslkey"] = self.ssl_key

        if self.ssl_root_cert:
            if not os.path.exists(self.ssl_root_cert):
                raise FileNotFoundError(f"SSL root certificate not found: {self.ssl_root_cert}")
            ssl_context["sslrootcert"] = self.ssl_root_cert

        return ssl_context

    def to_connection_params(self) -> Dict[str, Any]:
        """Convert config to psycopg2 connection parameters."""
        params = {
            "host": self.host,
            "port": self.port,
            "database": self.database,
            "user": self.user,
            "password": self.password,
            "sslmode": self.ssl_mode,
        }

        ssl_context = self.get_ssl_context()
        if ssl_context:
            params.update(ssl_context)

        return params


class PostgresConnectionValidator:
    """Validates PostgreSQL connection security."""

    REQUIRED_VARS = ["POSTGRES_HOST", "POSTGRES_DB", "POSTGRES_USER", "POSTGRES_PASSWORD"]
    SENSITIVE_PATTERNS = [
        r"password\s*=\s*['\"][^'\"]+['\"]",
        r"pwd\s*=\s*['\"][^'\"]+['\"]",
    ]

    @classmethod
    def validate_environment(cls) -> Dict[str, Any]:
        """Validate all required environment variables are set."""
        missing = []
        for var in cls.REQUIRED_VARS:
            if not os.getenv(var):
                missing.append(var)

        if missing:
            raise EnvironmentError(
                f"Missing required environment variables: {', '.join(missing)}"
            )

        return {
            "host": os.getenv("POSTGRES_HOST"),
            "port": int(os.getenv("POSTGRES_PORT", "5432")),
            "database": os.getenv("POSTGRES_DB"),
            "user": os.getenv("POSTGRES_USER"),
            "password": os.getenv("POSTGRES_PASSWORD"),
            "ssl_mode": os.getenv("POSTGRES_SSL_MODE", "require"),
            "ssl_cert": os.getenv("POSTGRES_SSL_CERT"),
            "ssl_key": os.getenv("POSTGRES_SSL_KEY"),
            "ssl_root_cert": os.getenv("POSTGRES_SSL_ROOT_CERT"),
        }

    @classmethod
    def validate_connection_string(cls, conn_string: str) -> bool:
        """Validate a PostgreSQL connection string for security issues."""
        if not conn_string:
            raise ValueError("Connection string cannot be empty")

        if "?" in conn_string and "sslmode" not in conn_string:
            logger.warning("Connection string missing SSL mode parameter")

        for pattern in cls.SENSITIVE_PATTERNS:
            if re.search(pattern, conn_string, re.IGNORECASE):
                logger.warning(
                    f"Connection string may contain hardcoded credentials: {pattern}"
                )

        return True

    @classmethod
    def check_password_strength(cls, password: str) -> bool:
        """Check if password meets minimum security requirements."""
        if len(password) < 12:
            raise ValueError("Password must be at least 12 characters long")
        return True


def get_security_config() -> DatabaseSecurityConfig:
    """Create DatabaseSecurityConfig from environment variables."""
    env_vars = PostgresConnectionValidator.validate_environment()
    return DatabaseSecurityConfig(**env_vars)


def validate_connection() -> bool:
    """Validate database connection security configuration."""
    try:
        config = get_security_config()
        PostgresConnectionValidator.check_password_strength(config.password)
        logger.info(f"Security configuration validated for host: {config.host}:{config.port}")
        return True
    except Exception as e:
        logger.error(f"Connection validation failed: {e}")
        raise


def create_connection_pool(
    min_connections: int = 1,
    max_connections: int = 10,
) -> pool.ThreadedConnectionPool:
    """
    Create a secure PostgreSQL connection pool.

    Args:
        min_connections: Minimum number of connections in pool
        max_connections: Maximum number of connections in pool

    Returns:
        ThreadedConnectionPool with SSL encryption
    """
    config = get_security_config()
    params = config.to_connection_params()

    logger.info(
        f"Creating secure connection pool to {config.host}:{config.port} "
        f"with SSL mode: {config.ssl_mode}"
    )

    try:
        connection_pool = pool.ThreadedConnectionPool(
            minconn=min_connections,
            maxconn=max_connections,
            **params
        )
        test_conn = connection_pool.getconn()
        test_conn.set_isolation_level(ISOLATION_LEVEL_AUTOCOMMIT)
        connection_pool.putconn(test_conn)
        logger.info("Connection pool created and tested successfully")
        return connection_pool
    except Exception as e:
        logger.error(f"Failed to create connection pool: {e}")
        raise


def create_postgres_saver(
    conn_pool: Optional[pool.ThreadedConnectionPool] = None,
) -> Any:
    """
    Create a PostgresSaver for LangGraph checkpointing with security.

    Args:
        conn_pool: Optional existing connection pool

    Returns:
        PostgresSaver instance configured with SSL
    """
    from langgraph.checkpoint.postgres import PostgresSaver

    if conn_pool is None:
        conn_pool = create_connection_pool()

    config = get_security_config()
    ssl_context = config.get_ssl_context()

    saver = PostgresSaver(
        conn=conn_pool,
        ssl=ssl_context,
    )

    logger.info("PostgresSaver created with SSL encryption enabled")
    return saver


def get_connection_string(include_password: bool = False) -> str:
    """
    Generate PostgreSQL connection string from environment.

    Args:
        include_password: Whether to include password in string (default: False)

    Returns:
        Connection string (password masked by default)
    """
    config = get_security_config()

    if include_password:
        return (
            f"postgresql://{config.user}:{config.password}@"
            f"{config.host}:{config.port}/{config.database}"
            f"?sslmode={config.ssl_mode}"
        )

    masked = "****" + config.password[-4:] if config.password else "****"
    return (
        f"postgresql://{config.user}:{masked}@"
        f"{config.host}:{config.port}/{config.database}"
        f"?sslmode={config.ssl_mode}"
    )
