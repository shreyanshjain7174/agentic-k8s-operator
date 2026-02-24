"""
LiteLLM client wrapper for multi-LLM routing.

Routes requests to gpt-4o (vision), gpt-4o-mini (text), or fallback models
via the in-cluster LiteLLM proxy. Never hardcodes API keys - they come from
Kubernetes Secrets mounted at /etc/secrets/llm-keys.

Cost tracking: gpt-4o = $0.01/analysis, gpt-4o-mini = $0.001/analysis
"""

import base64
import json
import logging
from typing import Optional, List, Dict
import os

import litellm
from litellm import completion

from agents.utils.credential_sanitizer import sanitize_credentials

logger = logging.getLogger(__name__)


class LiteLLMClientError(Exception):
    """Base exception for LiteLLM operations"""
    pass


class LiteLLMClient:
    """
    Wrapper for LiteLLM proxy client.
    
    Routes vision tasks (image analysis) to gpt-4o
    Routes text tasks to gpt-4o-mini (cheaper)
    """
    
    def __init__(
        self,
        proxy_url: Optional[str] = None,
        api_key_path: str = "/etc/secrets/llm-keys/openai-key",
        budget_limit_per_task: float = 0.02,
    ):
        """
        Initialize LiteLLM client.
        
        Args:
            proxy_url: LiteLLM proxy URL (default: inferred from env or localhost)
            api_key_path: Path to API key file (from Kubernetes Secret)
            budget_limit_per_task: Max spend per LLM call in USD
        """
        # Determine proxy URL
        self.proxy_url = proxy_url or os.getenv("LITELLM_PROXY_URL", "http://localhost:8000")
        self.api_key_path = api_key_path
        self.budget_limit = budget_limit_per_task
        
        # Load API key from Secret
        self.api_key = self._load_api_key()
        
        # Configure litellm
        litellm.api_base = self.proxy_url
        if self.api_key:
            litellm.api_key = self.api_key
        
        logger.info(f"LiteLLM client configured: proxy_url={self.proxy_url}")
    
    def _load_api_key(self) -> Optional[str]:
        """Load API key from Kubernetes Secret mount"""
        try:
            if os.path.exists(self.api_key_path):
                with open(self.api_key_path, "r") as f:
                    key = f.read().strip()
                    logger.info(f"Loaded API key from {self.api_key_path}")
                    return key
            else:
                logger.warning(f"API key not found at {self.api_key_path}, using proxy defaults")
                return None
        except Exception as e:
            logger.error(f"Error loading API key: {e}")
            return None
    
    async def analyze_screenshot(
        self,
        image_bytes: bytes,
        prompt: str,
        url: Optional[str] = None
    ) -> str:
        """
        Analyze screenshot using vision LLM (gpt-4o).
        
        Args:
            image_bytes: PNG screenshot bytes
            prompt: Analysis prompt (e.g., "What is the primary CTA?")
            url: Source URL (for logging)
            
        Returns:
            Analysis text
            
        Raises:
            LiteLLMClientError: On LLM errors
        """
        try:
            # Convert to base64 data URL
            image_b64 = base64.b64encode(image_bytes).decode("utf-8")
            image_url = f"data:image/png;base64,{image_b64}"
            
            logger.info(f"Analyzing screenshot via gpt-4o (url={url})")
            
            response = completion(
                model="gpt-4o",  # Vision model
                messages=[{
                    "role": "user",
                    "content": [
                        {
                            "type": "image_url",
                            "image_url": {"url": image_url}
                        },
                        {
                            "type": "text",
                            "text": prompt
                        }
                    ]
                }],
                max_tokens=500,
                temperature=0.7,
            )
            
            analysis = response.choices[0].message.content
            logger.info(f"Vision analysis complete for {url}: {len(analysis)} chars")
            return analysis
            
        except Exception as e:
            logger.error(f"Vision analysis failed: {e}")
            raise LiteLLMClientError(f"Vision analysis error: {str(e)}")
    
    async def synthesize_report(
        self,
        context: str,
        prompt: str
    ) -> str:
        """
        Synthesize competitive report using text LLM (gpt-4o-mini).
        
        Args:
            context: Compiled analysis context
            prompt: Synthesis prompt
            
        Returns:
            Generated report section
            
        Raises:
            LiteLLMClientError: On LLM errors
        """
        try:
            logger.info("Synthesizing report via gpt-4o-mini")
            
            response = completion(
                model="gpt-4o-mini",  # Cheaper text model
                messages=[{
                    "role": "user",
                    "content": f"{prompt}\n\nContext:\n{context}"
                }],
                max_tokens=2000,
                temperature=0.7,
            )
            
            report = response.choices[0].message.content
            logger.info(f"Report synthesis complete: {len(report)} chars")
            return report
            
        except Exception as e:
            logger.error(f"Report synthesis failed: {e}")
            raise LiteLLMClientError(f"Synthesis error: {str(e)}")
    
    async def extract_json(
        self,
        text: str,
        schema: str,
        model: str = "gpt-4o-mini"
    ) -> Dict:
        """
        Extract structured JSON from text.
        
        Args:
            text: Input text
            schema: JSON schema or description
            model: LLM to use
            
        Returns:
            Parsed JSON dict
            
        Raises:
            LiteLLMClientError: On parsing errors
        """
        try:
            prompt = f"""Extract JSON matching this schema from the text:

Schema: {schema}

Text: {text}

Return ONLY valid JSON, no markdown, no explanation."""
            
            response = completion(
                model=model,
                messages=[{"role": "user", "content": prompt}],
                temperature=0.0,  # Deterministic for JSON extraction
            )
            
            result_text = response.choices[0].message.content
            
            # Parse JSON
            result = json.loads(result_text)
            logger.info(f"JSON extraction successful: {len(result)} fields")
            return result
            
        except json.JSONDecodeError as e:
            logger.error(f"JSON parsing failed: {e}")
            raise LiteLLMClientError(f"JSON parsing error: {str(e)}")
        except Exception as e:
            logger.error(f"JSON extraction failed: {e}")
            raise LiteLLMClientError(f"Extraction error: {str(e)}")


# Singleton instance
_litellm_client: Optional[LiteLLMClient] = None


async def get_litellm_client() -> LiteLLMClient:
    """Get or create LiteLLM client (singleton)"""
    global _litellm_client
    if _litellm_client is None:
        _litellm_client = LiteLLMClient()
    return _litellm_client
