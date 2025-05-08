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
    const messageBox = document.getElementById('message');

    let ranges = [];

    // Format time as hh:mm:ss
    function formatTime(seconds) {
        const h = String(Math.floor(seconds / 3600)).padStart(2, '0');
        const m = String(Math.floor((seconds % 3600) / 60)).padStart(2, '0');
        const s = String(Math.floor(seconds % 60)).padStart(2, '0');
        return `${h}:${m}:${s}`;
    }

    // Fill start input from video player time
    setStartBtn.addEventListener('click', () => {
        startInput.value = formatTime(player.currentTime);
    });

    // Fill end input from video player time
    setEndBtn.addEventListener('click', () => {
        endInput.value = formatTime(player.currentTime);
    });

    // Show input block when adding a range
    showRangeInputsBtn.addEventListener('click', () => {
        rangeInputs.style.display = 'block';
        startInput.value = '';
        endInput.value = '';
    });

    // Save the current range and add it to the list
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

        // Remove range button
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

    // Submit the cut request to the backend
    cutForm.addEventListener('submit', async (e) => {
        messageBox.style.display = 'block';

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
                ranges
            })
        });

        const data = await res.json();

        messageBox.innerHTML = `
            <div class="d-flex justify-content-between align-items-center mt-4 mb-2">
                <p class="">✅ Video saved: <strong>${data.filename}</strong></p>
                <a class="btn btn-outline-success" href="/uploads/${data.filename}" download>Download</a>
            </div>
        `;
    });
});
