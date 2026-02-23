"""
Shared state schema for LangGraph workflow.

Defines the TypedDict that flows through all agent nodes.
Checkpointed by PostgresSaver for pod preemption resilience.
"""

from typing import TypedDict, Annotated, List, Optional, Dict
from langgraph.graph import add_messages


class AgentWorkflowState(TypedDict, total=False):
    """
    Shared state for visual scraping workflow.
    
    Used by all agent nodes. LangGraph checkpoints this state to PostgreSQL
    so that if a pod is preempted/killed mid-execution, the workflow resumes
    from the last completed node with state intact.
    """
    
    # ===== Input Parameters =====
    job_id: str
    """Unique identifier for this job (from AgentWorkload.metadata.name)"""
    
    target_urls: List[str]
    """List of competitor URLs to scrape and analyze"""
    
    # ===== Scraping Artifacts =====
    raw_html: Dict[str, str]
    """Mapping of URL -> raw HTML content. Populated by scraper_agent"""
    
    screenshots: Dict[str, bytes]
    """Mapping of URL -> PNG screenshot bytes. Populated by scraper_agent"""
    
    # ===== Analysis Results =====
    dom_structures: Dict[str, Dict]
    """
    Mapping of URL -> parsed DOM structure.
    Contains extracted pricing, CTAs, navigation, key elements.
    Populated by dom_agent.
    """
    
    visual_insights: List[Dict]
    """
    List of visual insights from screenshot analysis.
    Populated by screenshot_agent (vision LLM analysis).
    Each item: {"url": str, "insight": str, "confidence": float}
    """
    
    competitive_signals: List[Dict]
    """
    Synthesized competitive signals from combined analysis.
    Populated by synthesis_agent.
    """
    
    # ===== Output Artifacts =====
    report_path: Optional[str]
    """S3/MinIO path to generated report (PDF or HTML)"""
    
    report_content: Optional[str]
    """Generated report content (markdown or HTML)"""
    
    # ===== Status Tracking =====
    status: str
    """
    Job status: running | awaiting_review | complete | failed
    Initialized to 'running' by entrypoint.
    """
    
    error: Optional[str]
    """Error message if status == 'failed'"""
    
    # ===== Execution Metadata =====
    start_time: Optional[float]
    """Unix timestamp when job started"""
    
    end_time: Optional[float]
    """Unix timestamp when job completed"""
    
    messages: Annotated[List, add_messages]
    """Message history for agent reasoning (populated by LangGraph)"""
