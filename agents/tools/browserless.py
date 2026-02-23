"""
Browserless CDP (Chrome DevTools Protocol) client.

Connects to Browserless via WebSocket to:
- Capture full-page screenshots
- Extract DOM HTML
- Parse computed styles for key elements (pricing, CTA, nav)
- 30-second timeout with automatic session cleanup
"""

import asyncio
import base64
import json
import logging
from typing import Optional, Dict, Tuple
from urllib.parse import urlparse

import aiohttp
import websockets

logger = logging.getLogger(__name__)


class BrowserlessError(Exception):
    """Base exception for Browserless operations"""
    pass


class BrowserlessUnavailableError(BrowserlessError):
    """Raised when Browserless is unreachable. Triggers Argo pod retry."""
    pass


class BrowserlessTimeoutError(BrowserlessError):
    """Raised when operation exceeds timeout"""
    pass


class BrowserlessClient:
    """
    Async CDP client for Browserless.
    
    Usage:
        client = BrowserlessClient("ws://browserless:3000")
        html, screenshot = await client.scrape_url("https://example.com")
    """
    
    def __init__(self, browserless_url: str = "ws://browserless:3000", timeout: int = 30):
        """
        Initialize Browserless client.
        
        Args:
            browserless_url: WebSocket URL (e.g., "ws://browserless:3000")
            timeout: Session timeout in seconds (hard limit before release)
        """
        self.browserless_url = browserless_url
        self.timeout = timeout
        self.session: Optional[aiohttp.ClientSession] = None
    
    async def __aenter__(self):
        """Async context manager entry"""
        self.session = aiohttp.ClientSession()
        return self
    
    async def __aexit__(self, exc_type, exc_val, exc_tb):
        """Async context manager exit - cleanup session"""
        if self.session:
            await self.session.close()
    
    async def scrape_url(self, url: str) -> Tuple[str, bytes]:
        """
        Scrape a single URL.
        
        Args:
            url: URL to scrape
            
        Returns:
            Tuple of (html_content, screenshot_bytes)
            
        Raises:
            BrowserlessUnavailableError: If Browserless is unreachable
            BrowserlessTimeoutError: If operation exceeds timeout
            BrowserlessError: On other errors
        """
        try:
            return await asyncio.wait_for(
                self._scrape_with_cdp(url),
                timeout=self.timeout
            )
        except asyncio.TimeoutError:
            logger.error(f"Timeout scraping {url} after {self.timeout}s")
            raise BrowserlessTimeoutError(f"Timeout scraping {url}")
        except ConnectionRefusedError as e:
            logger.error(f"Cannot connect to Browserless: {e}")
            raise BrowserlessUnavailableError(f"Browserless unreachable at {self.browserless_url}")
        except Exception as e:
            logger.error(f"Error scraping {url}: {e}")
            raise BrowserlessError(f"Scraping failed for {url}: {str(e)}")
    
    async def _scrape_with_cdp(self, url: str) -> Tuple[str, bytes]:
        """
        Internal: Perform scraping via CDP WebSocket.
        
        Args:
            url: URL to scrape
            
        Returns:
            Tuple of (html_content, screenshot_bytes)
        """
        # Build Browserless CDP endpoint
        # Format: ws://browserless:3000/chromium/playwright?[options]
        browserless_endpoint = f"{self.browserless_url}/chromium/playwright"
        
        logger.info(f"Connecting to Browserless: {browserless_endpoint}")
        
        try:
            async with websockets.connect(browserless_endpoint) as websocket:
                logger.info(f"Connected to Browserless, navigating to {url}")
                
                # Navigate to URL
                await websocket.send(json.dumps({
                    "method": "Page.navigate",
                    "params": {"url": url}
                }))
                response = json.loads(await websocket.recv())
                
                if "error" in response:
                    raise BrowserlessError(f"CDP navigation failed: {response['error']}")
                
                # Wait for page load
                await asyncio.sleep(2)  # Give page time to render
                
                # Extract full HTML
                logger.info(f"Extracting HTML from {url}")
                await websocket.send(json.dumps({
                    "method": "Runtime.evaluate",
                    "params": {
                        "expression": "document.documentElement.outerHTML"
                    }
                }))
                html_response = json.loads(await websocket.recv())
                html_content = html_response.get("result", {}).get("value", "")
                
                if not html_content:
                    raise BrowserlessError("Failed to extract HTML")
                
                # Capture screenshot
                logger.info(f"Capturing screenshot from {url}")
                await websocket.send(json.dumps({
                    "method": "Page.captureScreenshot",
                    "params": {"format": "png"}
                }))
                screenshot_response = json.loads(await websocket.recv())
                screenshot_b64 = screenshot_response.get("result", {}).get("data", "")
                
                if not screenshot_b64:
                    raise BrowserlessError("Failed to capture screenshot")
                
                # Decode base64 screenshot
                screenshot_bytes = base64.b64decode(screenshot_b64)
                
                logger.info(f"Successfully scraped {url}: {len(html_content)} bytes HTML, {len(screenshot_bytes)} bytes screenshot")
                
                return html_content, screenshot_bytes
                
        except websockets.exceptions.WebSocketException as e:
            logger.error(f"WebSocket error: {e}")
            raise BrowserlessUnavailableError(f"WebSocket error: {str(e)}")
    
    async def extract_dom_structure(self, html_content: str) -> Dict:
        """
        Parse HTML and extract DOM structure.
        
        Extracts: pricing elements, CTAs, navigation, key sections
        
        Args:
            html_content: Raw HTML string
            
        Returns:
            Dict containing parsed structure
        """
        try:
            # Simple DOM extraction (can be enhanced with BeautifulSoup)
            structure = {
                "has_pricing": "pricing" in html_content.lower() or "$" in html_content,
                "has_cta": any(cta in html_content.lower() for cta in [
                    "button", "sign up", "get started", "try now", "contact"
                ]),
                "has_navigation": "<nav>" in html_content or "<header>" in html_content,
                "html_length": len(html_content),
                "title": self._extract_title(html_content),
            }
            return structure
        except Exception as e:
            logger.error(f"Error parsing DOM: {e}")
            return {"error": str(e)}
    
    def _extract_title(self, html_content: str) -> Optional[str]:
        """Extract page title from HTML"""
        try:
            start = html_content.find("<title>") + 7
            end = html_content.find("</title>")
            if start > 6 and end > start:
                return html_content[start:end]
        except Exception:
            pass
        return None


# Singleton instance
_browserless_client: Optional[BrowserlessClient] = None


async def get_browserless_client(url: str = "ws://browserless:3000") -> BrowserlessClient:
    """Get or create Browserless client (singleton)"""
    global _browserless_client
    if _browserless_client is None:
        _browserless_client = BrowserlessClient(url)
    return _browserless_client
