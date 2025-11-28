"""
Neo-Janus Red Agent: Automated Attack Campaign
Simulates adversarial prompts against Blue Sentinel defense
"""

import logging
import os
import random
import sys
import time
from typing import Dict, List, Optional

import requests
import yaml
from attack_lib.fuzzer import PromptFuzzer
from rich.console import Console
from rich.panel import Panel
from rich.progress import track

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(levelname)s] %(message)s"
)
logger = logging.getLogger(__name__)

console = Console()


class RedAgentAttacker:
    """Red Agent for adversarial attack simulation"""

    def __init__(self, config_path: str = "../config.yaml"):
        """Initialize attacker with configuration"""
        self.config = self._load_config(config_path)
        self.target_url = self.config.get("red_agent", {}).get("target_url")
        self.payloads = self._load_payloads()
        self.fuzzer = PromptFuzzer()
        self.session = requests.Session()
        self.session.timeout = 5

    def _load_config(self, config_path: str) -> Dict:
        """Load configuration from YAML file"""
        base_dir = os.path.dirname(os.path.abspath(__file__))
        full_path = os.path.join(base_dir, config_path)

        try:
            with open(full_path, "r", encoding="utf-8") as f:
                return yaml.safe_load(f)
        except FileNotFoundError:
            console.print(f"[red]Error:[/red] Config file not found: {full_path}")
            sys.exit(1)
        except yaml.YAMLError as e:
            console.print(f"[red]Error:[/red] Invalid YAML: {e}")
            sys.exit(1)

    def _load_payloads(self) -> List[str]:
        """Load base payloads from file"""
        base_dir = os.path.dirname(os.path.abspath(__file__))
        payloads_file = os.path.join(base_dir, "attack_lib/payloads.txt")

        try:
            with open(payloads_file, "r", encoding="utf-8") as f:
                payloads = [line.strip() for line in f if line.strip()]
                if not payloads:
                    logger.warning("Payloads file is empty")
                    return []
                logger.info(f"Loaded {len(payloads)} base payloads")
                return payloads
        except FileNotFoundError:
            console.print(f"[red]Error:[/red] Payloads file not found: {payloads_file}")
            return []

    def run_attack_campaign(self, intensity: int = 10) -> None:
        """Execute attack campaign"""
        if not self.payloads:
            console.print("[red]Error:[/red] No payloads available. Aborting.")
            return

        console.print(
            Panel.fit(
                f"🚀 INITIATING RED TEAM ATTACK CAMPAIGN\n"
                f"Target: {self.target_url}\n"
                f"Intensity: {intensity} requests",
                style="bold red"
            )
        )

        results = {"blocked": 0, "passed": 0, "errors": 0}

        for i in track(range(intensity), description="Firing payloads..."):
            base_payload = random.choice(self.payloads)
            fuzzed_payload = self.fuzzer.generate_payload(base_payload)

            try:
                status = self._send_attack(fuzzed_payload)
                results[status] += 1

                if status == "passed":
                    console.print(
                        f"[bold red]❌ VULNERABILITY FOUND:[/bold red] {fuzzed_payload[:50]}..."
                    )
                    logger.warning(f"Successful bypass detected: {fuzzed_payload[:100]}")

            except requests.RequestException as e:
                results["errors"] += 1
                logger.error(f"Attack request failed: {e}")

            # Random delay to avoid detection/overload
            time.sleep(random.uniform(0.05, 0.3))

        self._print_results(results, intensity)

    def _send_attack(self, payload: str) -> str:
        """Send attack payload and return status"""
        try:
            response = self.session.post(
                self.target_url,
                json={"input": payload, "source": "RED_AGENT"},
                timeout=5
            )

            if response.status_code == 200:
                data = response.json()
                status_key = data.get("status", "").upper()
                
                if status_key == "BLOCKED":
                    return "blocked"
                elif status_key == "APPROVED":
                    return "passed"
                else:
                    logger.warning(f"Unknown status: {status_key}")
                    return "errors"
            else:
                logger.error(f"HTTP {response.status_code}: {response.text[:100]}")
                return "errors"

        except requests.Timeout:
            logger.error("Request timeout")
            return "errors"
        except requests.ConnectionError:
            logger.error("Connection refused - is Janus Core running?")
            return "errors"
        except ValueError as e:
            logger.error(f"Invalid JSON response: {e}")
            return "errors"

    def _print_results(self, results: Dict[str, int], total: int) -> None:
        """Print campaign results"""
        bypass_rate = (results["passed"] / total * 100) if total > 0 else 0

        console.print(
            Panel(
                f"""[bold]ATTACK CAMPAIGN SUMMARY[/bold]
🛡️ Blocked by Blue Sentinel:  [green]{results['blocked']}[/green]
💀 Successful bypasses:         [red]{results['passed']}[/red] ({bypass_rate:.1f}%)
⚠️  Technical errors:           [yellow]{results['errors']}[/yellow]

[dim]Bypasses trigger Vaccine mechanism (threshold: {self.config.get('vaccine', {}).get('trigger_count', 'N/A')})[/dim]
""",
                title="Red Team Results",
                expand=False
            )
        )

        logger.info(
            f"Campaign complete: {results['blocked']} blocked, "
            f"{results['passed']} passed, {results['errors']} errors"
        )


def main() -> None:
    """Main entry point"""
    attacker = RedAgentAttacker()
    
    # Allow command-line intensity parameter
    intensity = int(sys.argv[1]) if len(sys.argv) > 1 else 10
    
    try:
        attacker.run_attack_campaign(intensity=intensity)
    except KeyboardInterrupt:
        console.print("\n[yellow]Attack campaign interrupted by user[/yellow]")
        sys.exit(0)
    except Exception as e:
        console.print(f"[red]Fatal error:[/red] {e}")
        logger.exception("Unexpected error during campaign")
        sys.exit(1)


if __name__ == "__main__":
    main()