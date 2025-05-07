document.addEventListener('DOMContentLoaded', () => {
    const cutForm = document.getElementById('cutForm');
    const startInput = document.getElementById('start');
    const endInput = document.getElementById('end');
    const deleteCheckbox = document.getElementById('deleteOriginal');
    const messageBox = document.getElementById('message');
    const setStartBtn = document.getElementById('setStart');
    const setEndBtn = document.getElementById('setEnd');
    const video = document.getElementById('player');

    if (!cutForm || !startInput || !endInput || !deleteCheckbox || !messageBox || !setStartBtn || !setEndBtn || !video) return;

    const formatTime = (sec) => {
        const h = String(Math.floor(sec / 3600)).padStart(2, '0');
        const m = String(Math.floor((sec % 3600) / 60)).padStart(2, '0');
        const s = String(Math.floor(sec % 60)).padStart(2, '0');
        return `${h}:${m}:${s}`;
    };

    setStartBtn.addEventListener('click', () => {
        const currentTime = Math.floor(video.currentTime);
        startInput.value = formatTime(currentTime);
    });

    setEndBtn.addEventListener('click', () => {
        const currentTime = Math.floor(video.currentTime);
        endInput.value = formatTime(currentTime);
    });

    cutForm.addEventListener('submit', async (e) => {
        e.preventDefault();

        const data = {
            filename: currentFilename,
            start: startInput.value,
            end: endInput.value,
            delete_original: deleteCheckbox.checked
        };

        try {
            const res = await fetch('/cut', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(data)
            });
            const result = await res.text();
            messageBox.innerText = `Saved: ${result}`;
        } catch (err) {
            messageBox.innerText = "Cut failed.";
        }
    });
});
