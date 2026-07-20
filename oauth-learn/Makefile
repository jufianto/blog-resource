.PHONY: all diagrams pdf lab-pdf lab2-pdf clean-preview

all: diagrams pdf lab-pdf lab2-pdf

diagrams:
	@mkdir -p resource/diagrams/exported
	@for f in resource/diagrams/*.mmd; do \
		name="$$(basename "$$f" .mmd)"; \
		npx -y @mermaid-js/mermaid-cli -i "$$f" -o "resource/diagrams/exported/$$name.png" -s 2; \
	done

pdf:
	@mkdir -p export/pdf
	@cd resource/lessons && for f in oauth_lesson_*.md; do \
		pandoc "$$f" --pdf-engine=tectonic --include-in-header=_header.tex \
			-V colorlinks=true -V urlcolor=teal -V linkcolor=teal \
			-o "../../export/pdf/$${f%.md}.pdf"; \
	done

lab-pdf:
	@mkdir -p export/labs
	@cd resource/labs/day1-go-oauth-demo && \
		pandoc LAB.md --pdf-engine=tectonic -o "../../../export/labs/day1_go_oauth_demo_lab.pdf"

lab2-pdf:
	@mkdir -p export/labs
	@cd resource/labs/day2-oauth-demo && \
		pandoc LAB.md --pdf-engine=tectonic -o "../../../export/labs/day2_oauth_advanced_lab.pdf"

clean-preview:
	@rm -rf export/preview
