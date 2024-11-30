let selectedCategories = [];


const categoryButtons = document.querySelectorAll('.category-btn');


function updateSelectedCategories() {
    
    document.getElementById('selectedCategories').value = selectedCategories.join(',');
}


categoryButtons.forEach(button => {
    button.addEventListener('click', function() {
        const category = this.getAttribute('data-category');
        
        
        if (selectedCategories.length >= 3 && !selectedCategories.includes(category)) {
            alert("Вы можете выбрать только 3 категории.");
            return; 
        }
        
        
        if (selectedCategories.includes(category)) {
            selectedCategories = selectedCategories.filter(item => item !== category);
            this.classList.remove('selected'); 
        } 
        
        else {
            selectedCategories.push(category);
            this.classList.add('selected'); 
        }

        
        updateSelectedCategories();
    });
});