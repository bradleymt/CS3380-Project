document.addEventListener('DOMContentLoaded', function() {
    
    // Library filter functionality
    const filterButtons = document.querySelectorAll('.filter-btn');
    filterButtons.forEach(button => {
        button.addEventListener('click', function() {
            // Remove active class from all buttons
            filterButtons.forEach(btn => btn.classList.remove('active'));
            // Add active class to clicked button
            this.classList.add('active');
            
            // Add filter the library items based on the selected filter
            console.log('Filter selected:', this.textContent);
        });
    });

    // Category buttons functionality
    const categoryButtons = document.querySelectorAll('.category-btn');
    categoryButtons.forEach(button => {
        button.addEventListener('click', function() {
            // Remove active class from all buttons
            categoryButtons.forEach(btn => btn.classList.remove('active'));
            // Add active class to clicked button
            this.classList.add('active');
            
            // Add filter the store games based on the selected category
            console.log('Category selected:', this.textContent);
        });
    });

    // Search functionality
    const searchInput = document.getElementById('searchInput');
    if (searchInput) {
        searchInput.addEventListener('input', function(e) {
            const searchTerm = e.target.value.toLowerCase();
            const gameCards = document.querySelectorAll('.game-card');
            
            gameCards.forEach(card => {
                const title = card.querySelector('.game-title').textContent.toLowerCase();
                if (title.includes(searchTerm)) {
                    card.style.display = 'block';
                } else {
                    card.style.display = 'none';
                }
            });
        });
    }

    // Add to cart functionality
    const addToCartButtons = document.querySelectorAll('.game-card .btn');
    addToCartButtons.forEach(button => {
        button.addEventListener('click', function(e) {
            e.stopPropagation();
            const gameCard = this.closest('.game-card');
            const gameTitle = gameCard.querySelector('.game-title').textContent;
            
            // Create a simple notification
            showNotification(`${gameTitle} added to cart!`);
        });
    });

    // Play button functionality in library
    const playButtons = document.querySelectorAll('.library-actions .btn');
    playButtons.forEach(button => {
        button.addEventListener('click', function() {
            const libraryItem = this.closest('.library-item');
            const gameTitle = libraryItem.querySelector('.library-title').textContent;
            
            showNotification(`Launching ${gameTitle}...`);
        });
    });

    // Game card click functionality
    const gameCards = document.querySelectorAll('.game-card');
    gameCards.forEach(card => {
        card.addEventListener('click', function() {
            const gameTitle = this.querySelector('.game-title').textContent;
            console.log('Game clicked:', gameTitle);
            // Add ability to open detailed game page or modal
        });
    });

    // Library item click functionality
    const libraryItems = document.querySelectorAll('.library-item');
    libraryItems.forEach(item => {
        item.addEventListener('click', function(e) {
            // Don't trigger if clicking on a button
            if (!e.target.closest('.btn')) {
                const gameTitle = this.querySelector('.library-title').textContent;
                console.log('Library game clicked:', gameTitle);
                // Add ability to open game details
            }
        });
    });

    // Notification function
    function showNotification(message) {
        // Create notification element
        const notification = document.createElement('div');
        notification.textContent = message;
        notification.style.cssText = `
            position: fixed;
            top: 80px;
            right: 20px;
            background: #66c0f4;
            color: #171a21;
            padding: 15px 25px;
            border-radius: 6px;
            font-weight: bold;
            z-index: 1000;
            animation: slideIn 0.3s ease-out;
            box-shadow: 0 4px 15px rgba(0,0,0,0.3);
        `;
        
        document.body.appendChild(notification);
        
        // Remove after 3 seconds
        setTimeout(() => {
            notification.style.animation = 'slideOut 0.3s ease-out';
            setTimeout(() => {
                document.body.removeChild(notification);
            }, 300);
        }, 3000);
    }

    // Add animation styles
    const style = document.createElement('style');
    style.textContent = `
        @keyframes slideIn {
            from {
                transform: translateX(400px);
                opacity: 0;
            }
            to {
                transform: translateX(0);
                opacity: 1;
            }
        }
        
        @keyframes slideOut {
            from {
                transform: translateX(0);
                opacity: 1;
            }
            to {
                transform: translateX(400px);
                opacity: 0;
            }
        }
    `;
    document.head.appendChild(style);

    // Smooth scroll for navigation
    document.querySelectorAll('a[href^="#"]').forEach(anchor => {
        anchor.addEventListener('click', function (e) {
            e.preventDefault();
            const target = document.querySelector(this.getAttribute('href'));
            if (target) {
                target.scrollIntoView({
                    behavior: 'smooth'
                });
            }
        });
    });

    console.log('GameVault initialized successfully!');
});