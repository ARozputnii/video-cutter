document.addEventListener('DOMContentLoaded', () => {
    const player = document.getElementById('player');
    const startInput = document.getElementById('start');
    const endInput = document.getElementById('end');
    const setStartBtn = document.getElementById('setStart');
    const setEndBtn = document.getElementById('setEnd');
    const showRangeInputsBtn = document.getElementById('showRangeInputs');
    const saveRangeBtn = document.getElementById('saveRange');
    const rangeInputs = document.getElementById('rangeInputs');
    const rangeList = document.getElementById('rangeList');
    const cutForm = document.getElementById('cutForm');
    const deleteCheckbox = document.getElementById('deleteOriginal');
    const messageBox = document.getElementById('message');

    let ranges = [];

    function formatTime(seconds) {
        const h = String(Math.floor(seconds / 3600)).padStart(2, '0');
        const m = String(Math.floor((seconds % 3600) / 60)).padStart(2, '0');
        const s = String(Math.floor(seconds % 60)).padStart(2, '0');
        return `${h}:${m}:${s}`;
    }

    setStartBtn.addEventListener('click', () => {
        startInput.value = formatTime(player.currentTime);
    });

    setEndBtn.addEventListener('click', () => {
        endInput.value = formatTime(player.currentTime);
    });

    showRangeInputsBtn.addEventListener('click', () => {
        rangeInputs.style.display = 'block';
        startInput.value = '';
        endInput.value = '';
    });

    saveRangeBtn.addEventListener('click', () => {
        const start = startInput.value;
        const end = endInput.value;
        if (!start || !end) {
            alert("Both start and end must be filled.");
            return;
        }

        ranges.push({ start, end });

        const li = document.createElement('li');
        li.className = 'list-group-item d-flex justify-content-between';
        li.textContent = `${start} → ${end}`;

        const removeBtn = document.createElement('button');
        removeBtn.className = 'btn btn-sm btn-danger';
        removeBtn.textContent = '×';
        removeBtn.onclick = () => {
            rangeList.removeChild(li);
            ranges = ranges.filter(r => !(r.start === start && r.end === end));
        };

        li.appendChild(removeBtn);
        rangeList.appendChild(li);

        rangeInputs.style.display = 'none';
    });

    cutForm.addEventListener('submit', async (e) => {
        e.preventDefault();

        if (ranges.length === 0) {
            alert("Please add at least one range.");
            return;
        }

        const res = await fetch('/cut', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                filename: currentFilename,
                ranges,
                delete_original: deleteCheckbox.checked
            })
        });

        const text = await res.text();
        messageBox.textContent = text;
    });
});
