async function deleteTopic() {
    const currentUrl = window.location.href; // Получение текущего URL

    if (confirm('Are you sure you want to delete this topic?')) {
        try {
            const response = await fetch(currentUrl, {
                method: 'DELETE',
            });

            if (response.ok) {
                alert('Topic deleted successfully!');
                // Перенаправление на профиль пользователя после удаления
                window.location.href = '/private/profile';
            } else {
                alert('Failed to delete the topic. Please try again.');
            }
        } catch (error) {
            console.error('Error deleting topic:', error);
            alert('An error occurred while deleting the topic.');
        }
    }
}