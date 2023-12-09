let shouldCloseApp = true

document.addEventListener("click", (e) => {
    if (e.target instanceof HTMLAnchorElement || e.target instanceof HTMLInputElement) {
        console.log(e.target.tagName);
        shouldCloseApp = false
    }
})

document.onvisibilitychange = async () => {
    if (shouldCloseApp) {
        await fetch("/shutdown")
    }
}