let currentFilename = "";

document.addEventListener('DOMContentLoaded', () => {
    const uploadForm = document.getElementById('uploadForm');
    const videoInput = document.getElementById('video');
    const videoPreview = document.getElementById('video-preview');
    const videoElement = document.getElementById('player');
    const cutFormWrapper = document.getElementById('cutFormWrapper');
    const messageBox = document.getElementById('message');
    const loader = document.getElementById('loader');

    if (!uploadForm || !videoInput || !videoPreview || !videoElement || !cutFormWrapper || !messageBox || !loader) return;

    uploadForm.addEventListener('submit', async (e) => {
        e.preventDefault();

        const videoFile = videoInput.files[0];
        if (!videoFile) {
            alert("Please select a video file.");
            return;
        }

        const formData = new FormData();
        formData.append('video', videoFile);

        loader.style.display = 'block';
        messageBox.style.display = 'none';

        try {
            const res = await fetch('/upload', {
                method: 'POST',
                body: formData
            });

            if (!res.ok) {
                throw new Error("Upload failed: " + res.statusText);
            }

            const filename = await res.text();

            currentFilename = filename;
            videoElement.src = `/uploads/${filename}`;
            videoPreview.style.display = 'block';
            cutFormWrapper.style.display = 'block';
            messageBox.innerText = "";
        } catch (err) {
            messageBox.innerText = "❌ Upload failed. " + err.message;
            messageBox.style.display = 'block';
        } finally {
            loader.style.display = 'none';
        }
    });
});
