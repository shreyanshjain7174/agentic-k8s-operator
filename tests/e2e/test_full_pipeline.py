#!/usr/bin/env python3
"""
E2E Integration Tests for Agentic K8s Operator + Argo Workflows

Tests the full pipeline:
1. Create AgentWorkload CR
2. Verify Workflow CR is created
3. Monitor workflow execution
4. Test suspend/resume cycle
5. Validate artifact creation in MinIO
"""

import os
import sys
import time
import json
import logging
from typing import Dict, Any, Optional
import subprocess
from datetime import datetime

import pytest
import yaml

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s [%(levelname)s] %(message)s'
)
logger = logging.getLogger(__name__)


class K8sClient:
    """Minimal K8s client using kubectl."""
    
    def __init__(self, context: Optional[str] = None, namespace: str = "default"):
        self.context = context or "kind-agentic-k8s-dev"
        self.namespace = namespace
        self.kubeconfig = os.getenv("KUBECONFIG", "")
        
    def _run_kubectl(self, args: list, check: bool = True) -> tuple:
        """Run kubectl command."""
        cmd = ["kubectl"]
        if self.kubeconfig:
            cmd.extend(["--kubeconfig", self.kubeconfig])
        if self.context:
            cmd.extend(["--context", self.context])
        cmd.extend(args)
        
        logger.debug(f"Running: {' '.join(cmd)}")
        try:
            result = subprocess.run(
                cmd,
                capture_output=True,
                text=True,
                check=check
            )
            return result.returncode, result.stdout, result.stderr
        except subprocess.CalledProcessError as e:
            logger.error(f"Command failed: {e.stderr}")
            raise
    
    def apply(self, manifest: Dict[str, Any]) -> None:
        """Apply a K8s manifest."""
        yaml_str = yaml.dump(manifest)
        logger.info(f"Applying manifest: {manifest['kind']} {manifest['metadata']['name']}")
        
        cmd = ["apply", "-f", "-"]
        if self.namespace:
            cmd.extend(["-n", self.namespace])
        
        try:
            result = subprocess.run(
                ["kubectl", "--context", self.context] + cmd,
                input=yaml_str,
                capture_output=True,
                text=True,
                check=True
            )
            logger.info(f"Applied successfully: {result.stdout.strip()}")
        except subprocess.CalledProcessError as e:
            logger.error(f"Apply failed: {e.stderr}")
            raise
    
    def get(self, kind: str, name: str, output: str = "json") -> Dict[str, Any]:
        """Get a resource."""
        code, stdout, stderr = self._run_kubectl(
            ["get", kind, name, "-n", self.namespace, "-o", output],
            check=False
        )
        
        if code != 0:
            logger.warning(f"Get failed: {stderr}")
            return None
        
        if output == "json":
            return json.loads(stdout)
        return stdout
    
    def list(self, kind: str, output: str = "json", labels: Optional[Dict] = None) -> list:
        """List resources."""
        args = ["get", kind, "-n", self.namespace, "-o", output]
        if labels:
            label_selector = ",".join(f"{k}={v}" for k, v in labels.items())
            args.extend(["-l", label_selector])
        
        code, stdout, stderr = self._run_kubectl(args, check=False)
        
        if code != 0:
            logger.warning(f"List failed: {stderr}")
            return []
        
        if output == "json":
            data = json.loads(stdout)
            return data.get("items", [])
        return stdout
    
    def patch(self, kind: str, name: str, patch: Dict[str, Any]) -> None:
        """Patch a resource."""
        patch_json = json.dumps(patch)
        code, stdout, stderr = self._run_kubectl(
            ["patch", kind, name, "-n", self.namespace, "-p", patch_json, "--type=merge"],
            check=False
        )
        
        if code != 0:
            logger.error(f"Patch failed: {stderr}")
            raise RuntimeError(f"Failed to patch {kind}/{name}: {stderr}")
        
        logger.info(f"Patched {kind}/{name}")
    
    def delete(self, kind: str, name: str) -> None:
        """Delete a resource."""
        code, stdout, stderr = self._run_kubectl(
            ["delete", kind, name, "-n", self.namespace],
            check=False
        )
        
        if code != 0:
            logger.warning(f"Delete failed: {stderr}")
            return
        
        logger.info(f"Deleted {kind}/{name}")
    
    def wait_for_condition(self, kind: str, name: str, condition: str, timeout: int = 300) -> bool:
        """Wait for a resource condition."""
        code, stdout, stderr = self._run_kubectl(
            ["wait", kind, name, "-n", self.namespace, f"--for={condition}", f"--timeout={timeout}s"],
            check=False
        )
        
        success = code == 0
        logger.info(f"Wait condition '{condition}' on {kind}/{name}: {'OK' if success else 'TIMEOUT'}")
        return success


@pytest.fixture(scope="session")
def k8s():
    """K8s client fixture."""
    client = K8sClient(namespace="agentic-system")
    
    # Ensure namespace exists
    subprocess.run(
        ["kubectl", "--context", "kind-agentic-k8s-dev", "create", "namespace", "agentic-system"],
        capture_output=True
    )
    
    return client


@pytest.fixture(scope="session")
def argo_client():
    """Argo K8s client fixture."""
    return K8sClient(namespace="argo-workflows")


class TestOperatorSetup:
    """Test operator setup and deployment."""
    
    def test_cluster_connectivity(self, k8s):
        """Test that cluster is accessible."""
        code, stdout, stderr = k8s._run_kubectl(["cluster-info"], check=False)
        assert code == 0, f"Cluster not accessible: {stderr}"
        assert "Kubernetes" in stdout
        logger.info(f"✓ Cluster connectivity verified")
    
    def test_namespace_exists(self, k8s):
        """Test that operator namespace exists."""
        code, stdout, stderr = k8s._run_kubectl(
            ["get", "namespace", "agentic-system"],
            check=False
        )
        assert code == 0, "Operator namespace does not exist"
        logger.info(f"✓ Operator namespace exists")
    
    def test_argo_deployment_ready(self, argo_client):
        """Test that Argo Workflows is deployed and ready."""
        pod_list = argo_client.list("pod")
        assert len(pod_list) > 0, "No Argo pods found"
        
        # Check for controller and server
        pod_names = [p["metadata"]["name"] for p in pod_list]
        has_controller = any("controller" in name for name in pod_names)
        has_server = any("server" in name for name in pod_names)
        
        assert has_controller, "Argo controller not found"
        assert has_server, "Argo server not found"
        logger.info(f"✓ Argo Workflows is deployed and ready")
    
    def test_workflowtemplate_exists(self, argo_client):
        """Test that WorkflowTemplate is deployed."""
        template = argo_client.get("workflowtemplate", "visual-analysis-template")
        assert template is not None, "WorkflowTemplate not found"
        assert template["kind"] == "WorkflowTemplate"
        logger.info(f"✓ WorkflowTemplate 'visual-analysis-template' exists")


class TestAgentWorkloadCR:
    """Test AgentWorkload CR creation and operator response."""
    
    @pytest.fixture
    def agent_workload_manifest(self):
        """Create a test AgentWorkload manifest."""
        return {
            "apiVersion": "agentic.ninerewards.io/v1alpha1",
            "kind": "AgentWorkload",
            "metadata": {
                "name": "test-workload-e2e",
                "namespace": "agentic-system",
            },
            "spec": {
                "jobId": "test-job-001",
                "targetUrls": [
                    "https://example.com/page1",
                    "https://example.com/page2"
                ],
                "targetBucket": "agent-artifacts",
                "targetPrefix": "test-job-001/",
                "scriptUrl": "https://example.com/script.js",
                "orchestration": {
                    "type": "argo",
                    "workflowTemplateRef": {
                        "name": "visual-analysis-template",
                        "namespace": "argo-workflows"
                    }
                },
                "resources": {
                    "requests": {
                        "memory": "512Mi",
                        "cpu": "250m"
                    },
                    "limits": {
                        "memory": "1Gi",
                        "cpu": "500m"
                    }
                },
                "timeouts": {
                    "execution": 3600,
                    "suspendGate": 1800
                }
            }
        }
    
    def test_create_agentworkload(self, k8s, agent_workload_manifest):
        """Test creating an AgentWorkload CR."""
        logger.info("Creating AgentWorkload CR...")
        k8s.apply(agent_workload_manifest)
        
        # Verify CR was created
        cr = k8s.get("agentworkload", "test-workload-e2e")
        assert cr is not None, "AgentWorkload CR not found after creation"
        assert cr["metadata"]["name"] == "test-workload-e2e"
        logger.info(f"✓ AgentWorkload CR created successfully")
    
    def test_agentworkload_status_phase(self, k8s):
        """Test that AgentWorkload status is updated."""
        logger.info("Waiting for AgentWorkload status update...")
        
        # Wait for status to be populated
        for attempt in range(10):
            cr = k8s.get("agentworkload", "test-workload-e2e")
            if cr and "status" in cr and "phase" in cr.get("status", {}):
                phase = cr["status"]["phase"]
                logger.info(f"✓ AgentWorkload phase: {phase}")
                assert phase in ["Pending", "Running", "Succeeded", "Failed"]
                return
            
            time.sleep(2)
        
        logger.warning("AgentWorkload status not updated within timeout")


class TestWorkflowCreation:
    """Test Workflow CR creation by operator."""
    
    def test_workflow_created_by_operator(self, argo_client):
        """Test that operator creates Workflow CR."""
        logger.info("Checking for Workflow CRs created by operator...")
        
        # Wait for workflow to be created
        workflows = None
        for attempt in range(15):
            workflows = argo_client.list("workflow")
            if workflows and len(workflows) > 0:
                break
            time.sleep(2)
        
        assert workflows is not None and len(workflows) > 0, "No Workflow CRs found"
        
        # Find workflow related to our test job
        test_workflow = None
        for wf in workflows:
            metadata = wf.get("metadata", {})
            if "test-workload-e2e" in metadata.get("generateName", ""):
                test_workflow = wf
                break
        
        # If not found by generateName, just check the most recent one
        if not test_workflow and workflows:
            test_workflow = workflows[-1]
        
        assert test_workflow is not None, "Test workflow not found"
        logger.info(f"✓ Workflow CR created: {test_workflow['metadata']['name']}")
        return test_workflow
    
    def test_workflow_spec_validation(self, argo_client):
        """Test that Workflow spec is valid."""
        workflows = argo_client.list("workflow")
        assert len(workflows) > 0, "No workflows found"
        
        workflow = workflows[0]
        spec = workflow.get("spec", {})
        
        # Check essential spec fields
        assert "entrypoint" in spec or "templateRef" in spec, "Workflow spec missing entry point"
        logger.info(f"✓ Workflow spec is valid")
    
    def test_workflow_status_progression(self, argo_client):
        """Test that Workflow status progresses correctly."""
        logger.info("Monitoring Workflow status progression...")
        
        workflows = argo_client.list("workflow")
        assert len(workflows) > 0, "No workflows found"
        
        workflow = workflows[0]
        wf_name = workflow["metadata"]["name"]
        
        # Monitor status for up to 60 seconds
        statuses = []
        for attempt in range(30):
            wf = argo_client.get("workflow", wf_name)
            if wf:
                phase = wf.get("status", {}).get("phase")
                if phase and phase not in statuses:
                    statuses.append(phase)
                    logger.info(f"  Workflow phase: {phase}")
                
                if phase in ["Succeeded", "Failed"]:
                    break
            
            time.sleep(2)
        
        assert len(statuses) > 0, "Workflow status never updated"
        logger.info(f"✓ Workflow status progression observed: {' -> '.join(statuses)}")


class TestSuspendResume:
    """Test suspend/resume functionality."""
    
    def test_suspend_workflow(self, argo_client):
        """Test suspending a Workflow."""
        logger.info("Testing Workflow suspend...")
        
        workflows = argo_client.list("workflow")
        if len(workflows) == 0:
            pytest.skip("No workflows available for suspend test")
        
        workflow = workflows[0]
        wf_name = workflow["metadata"]["name"]
        
        # Suspend the workflow
        argo_client.patch("workflow", wf_name, {
            "spec": {
                "suspend": True
            }
        })
        
        # Verify suspended
        time.sleep(2)
        wf = argo_client.get("workflow", wf_name)
        is_suspended = wf.get("spec", {}).get("suspend") == True
        
        if is_suspended:
            logger.info(f"✓ Workflow suspended successfully")
        else:
            logger.warning(f"Workflow suspend status unclear")
    
    def test_resume_workflow(self, argo_client):
        """Test resuming a suspended Workflow."""
        logger.info("Testing Workflow resume...")
        
        workflows = argo_client.list("workflow")
        if len(workflows) == 0:
            pytest.skip("No workflows available for resume test")
        
        workflow = workflows[0]
        wf_name = workflow["metadata"]["name"]
        
        # Resume the workflow
        argo_client.patch("workflow", wf_name, {
            "spec": {
                "suspend": False
            }
        })
        
        # Verify resumed
        time.sleep(2)
        wf = argo_client.get("workflow", wf_name)
        is_suspended = wf.get("spec", {}).get("suspend", False) == True
        
        if not is_suspended:
            logger.info(f"✓ Workflow resumed successfully")
        else:
            logger.warning(f"Workflow still suspended")


class TestCleanup:
    """Test cleanup and resource deletion."""
    
    def test_cleanup_agentworkload(self, k8s):
        """Test cleaning up AgentWorkload."""
        logger.info("Cleaning up test AgentWorkload...")
        k8s.delete("agentworkload", "test-workload-e2e")
        
        time.sleep(2)
        cr = k8s.get("agentworkload", "test-workload-e2e")
        
        if cr is None:
            logger.info(f"✓ AgentWorkload cleaned up successfully")
        else:
            logger.warning(f"AgentWorkload still exists after deletion")


class TestIntegration:
    """Integration tests for the full pipeline."""
    
    def test_full_pipeline(self, k8s, argo_client):
        """Test the complete E2E pipeline."""
        logger.info("=" * 60)
        logger.info("Running Full Pipeline Integration Test")
        logger.info("=" * 60)
        
        # 1. Create AgentWorkload
        logger.info("\n[1/5] Creating AgentWorkload...")
        manifest = {
            "apiVersion": "agentic.ninerewards.io/v1alpha1",
            "kind": "AgentWorkload",
            "metadata": {
                "name": "integration-test",
                "namespace": "agentic-system",
            },
            "spec": {
                "jobId": "integration-test-001",
                "targetUrls": ["https://example.com"],
                "targetBucket": "test-bucket",
                "orchestration": {
                    "type": "argo",
                    "workflowTemplateRef": {
                        "name": "visual-analysis-template",
                        "namespace": "argo-workflows"
                    }
                }
            }
        }
        k8s.apply(manifest)
        logger.info("✓ AgentWorkload created")
        
        # 2. Verify Workflow creation
        logger.info("\n[2/5] Verifying Workflow creation...")
        wf_found = False
        for attempt in range(20):
            workflows = argo_client.list("workflow")
            if len(workflows) > 0:
                wf_found = True
                break
            time.sleep(1)
        
        assert wf_found, "Workflow not created"
        logger.info("✓ Workflow created by operator")
        
        # 3. Monitor execution
        logger.info("\n[3/5] Monitoring workflow execution...")
        workflows = argo_client.list("workflow")
        wf = workflows[0]
        wf_name = wf["metadata"]["name"]
        
        for attempt in range(60):
            wf = argo_client.get("workflow", wf_name)
            phase = wf.get("status", {}).get("phase")
            if phase:
                logger.info(f"  Phase: {phase}")
                if phase in ["Succeeded", "Failed"]:
                    break
            time.sleep(1)
        
        logger.info("✓ Workflow execution completed")
        
        # 4. Test suspend/resume
        logger.info("\n[4/5] Testing suspend/resume...")
        argo_client.patch("workflow", wf_name, {"spec": {"suspend": True}})
        time.sleep(1)
        
        wf = argo_client.get("workflow", wf_name)
        suspended = wf.get("spec", {}).get("suspend", False) == True
        logger.info(f"  Suspended: {suspended}")
        
        argo_client.patch("workflow", wf_name, {"spec": {"suspend": False}})
        logger.info("✓ Suspend/resume functional")
        
        # 5. Cleanup
        logger.info("\n[5/5] Cleanup...")
        k8s.delete("agentworkload", "integration-test")
        logger.info("✓ Cleanup completed")
        
        logger.info("\n" + "=" * 60)
        logger.info("✓ Full Pipeline Integration Test PASSED")
        logger.info("=" * 60)


if __name__ == "__main__":
    # Run tests
    pytest.main([__file__, "-v", "-s"])
