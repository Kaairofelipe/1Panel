#!/bin/bash
cat << 'INNER_EOF' > patch.diff
<<<<<<< SEARCH
      - uses: fit2cloud/LLM-CodeReview-Action@main
=======
      - uses: fit2cloud/LLM-CodeReview-Action@016b651a0e3c47994f2ae630415986c89694d01a
>>>>>>> REPLACE
INNER_EOF
patch .github/workflows/llm-code-review.yml < patch.diff
