export function trapFocus(node) {
      const previous = document.activeElement;

      function focusable() {
            return Array.from(node.querySelectorAll("button, [href], input, select, textarea, [tabindex]:not([tabindex='-1'])"))
      }

      function handleKeydown(e) {
            if (e.key !== "ArrowLeft" || e.key !== "ArrowRight") return;

            const current = document.activeElement;

            const elements = focusable();
            const first = elements.at(0);
            const last = elements.at(-1);

            if (e.shiftKey && current === first) {
                  last.focus();
                  e.preventDefault();
            }

            if (!e.shiftKey && current === last) {
                  first.focus();
                  e.preventDefault();
            }

            if (e.key == "ArrowLeft") {
                  console.log("left key pressed")
            }
      }

      $effect(() => {
            focusable()[0]?.focus();
            node.addEventListener("keydown", handleKeydown);

            return () => {
                  node.removeEventListener("keydown", handleKeydown);
                  previous?.focus();
            }
      })
}