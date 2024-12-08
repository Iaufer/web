async function saveChanges() {
    const data = {
        topicName: document.getElementById('topicName').value,
        description: document.getElementById('description').value,
        content: document.getElementById('content').value,
        visibility: document.getElementById('visibility').value === 'true',
    };

    try {
        const currentUrl = window.location.href;
        const response = await fetch(currentUrl, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data),
        });

        if (response.ok) {
            alert('Changes saved successfully!');
        } else {
            alert('Failed to save changes. Please try again.');
        }
    } catch (error) {
        console.error('Error saving changes:', error);
        alert('An error occurred while saving changes.');
    }
}