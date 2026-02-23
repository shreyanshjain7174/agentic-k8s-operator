"""
Tests for LangGraph state schema.
"""

import pytest
from agents.graph.state import AgentWorkflowState


def test_state_schema_structure():
    """Test that state schema has all required fields"""
    state: AgentWorkflowState = {
        "job_id": "test-001",
        "target_urls": ["https://example.com"],
        "status": "running",
        "raw_html": {"https://example.com": "<html></html>"},
        "screenshots": {"https://example.com": b"data"},
        "dom_structures": {},
        "visual_insights": [],
        "competitive_signals": [],
        "report_content": None,
        "error": None,
        "start_time": None,
        "end_time": None,
        "messages": [],
    }
    
    assert state["job_id"] == "test-001"
    assert len(state["target_urls"]) == 1
    assert state["status"] == "running"
    print("✅ State schema valid")


def test_state_initialization():
    """Test state can be created with defaults"""
    state: AgentWorkflowState = {
        "job_id": "test",
        "target_urls": [],
        "status": "running",
        "raw_html": {},
        "screenshots": {},
        "dom_structures": {},
        "visual_insights": [],
        "competitive_signals": [],
        "report_content": None,
        "error": None,
        "start_time": None,
        "end_time": None,
        "messages": [],
    }
    
    # All required fields present
    assert "job_id" in state
    assert "target_urls" in state
    assert "status" in state
    print("✅ State initialization works")


def test_state_mutations():
    """Test that state can be mutated"""
    state: AgentWorkflowState = {
        "job_id": "test",
        "target_urls": ["https://example.com"],
        "status": "running",
        "raw_html": {},
        "screenshots": {},
        "dom_structures": {},
        "visual_insights": [],
        "competitive_signals": [],
        "report_content": None,
        "error": None,
        "start_time": None,
        "end_time": None,
        "messages": [],
    }
    
    # Mutate state
    state["status"] = "complete"
    state["report_content"] = "Generated report"
    state["raw_html"]["https://example.com"] = "<html>Test</html>"
    
    assert state["status"] == "complete"
    assert state["report_content"] == "Generated report"
    assert len(state["raw_html"]) > 0
    print("✅ State mutations work")


def test_state_error_handling():
    """Test state error field"""
    state: AgentWorkflowState = {
        "job_id": "error-test",
        "target_urls": [],
        "status": "failed",
        "raw_html": {},
        "screenshots": {},
        "dom_structures": {},
        "visual_insights": [],
        "competitive_signals": [],
        "report_content": None,
        "error": "Connection refused",
        "start_time": None,
        "end_time": None,
        "messages": [],
    }
    
    assert state["status"] == "failed"
    assert state["error"] is not None
    assert "Connection" in state["error"]
    print("✅ Error handling in state works")
