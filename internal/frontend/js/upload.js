let currentFilename = "";

document.addEventListener('DOMContentLoaded', () => {
    const uploadForm = document.getElementById('uploadForm');
    const videoInput = document.getElementById('video');
    const videoPreview = document.getElementById('video-preview');
    const videoElement = document.getElementById('player');
    const cutForm = document.getElementById('cutFormWrapper');
    const messageBox = document.getElementById('message');

    if (!uploadForm || !videoInput || !videoPreview || !videoElement || !cutForm || !messageBox) return;

    uploadForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const formData = new FormData();
        const videoFile = videoInput.files[0];
        formData.append('video', videoFile);

        try {
            const res = await fetch('/upload', { method: 'POST', body: formData });
            const filename = await res.text();

            currentFilename = filename;
            videoElement.src = `/uploads/${filename}`;
            videoPreview.style.display = 'block';
            cutForm.style.display = 'block';
            messageBox.innerText = "";
        } catch (err) {
            messageBox.innerText = "Upload failed.";
        }
    });
});
