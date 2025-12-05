const API_BASE = '/api';
const SECURE_BASE = '/api/secure';

let cart = [];
let user = null;
let currentCategory = 'All';
let authToken = localStorage.getItem('authToken');

// Wrapper for all API calls to handle auth headers and errors
async function fetchAPI(endpoint, options = {}, requiresAuth = false) {
    try {
        const url = requiresAuth ? `${SECURE_BASE}${endpoint}` : `${API_BASE}${endpoint}`;
        
        if (authToken && requiresAuth) {
            options.headers = {
                ...options.headers,
                'Authorization': `Bearer ${authToken}`
            };
        }
        
        const response = await fetch(url, options);
        if (!response.ok) {
            const error = await response.text();
            if (response.status === 401) {
                localStorage.removeItem('authToken');
                localStorage.removeItem('username');
                authToken = null;
                window.location.href = 'login.html';
                throw new Error('Session expired');
            }
            throw new Error(error || 'API request failed');
        }
        return await response.json();
    } catch (error) {
        console.error('API Error:', error);
        if (!error.message.includes('Session expired')) {
            showNotification(error.message || 'An error occurred', 'error');
        }
        throw error;
    }
}

function formatPrice(price) {
    return `$${price.toFixed(2)}`;
}

// Show toast notification in top right corner
function showNotification(message, type = 'success') {
    const notification = document.createElement('div');
    notification.textContent = message;
    notification.className = `notification ${type}`;
    notification.style.cssText = `
        position: fixed;
        top: 80px;
        right: 20px;
        background: ${type === 'error' ? '#c74646' : '#66c0f4'};
        color: ${type === 'error' ? '#fff' : '#171a21'};
        padding: 15px 25px;
        border-radius: 6px;
        font-weight: bold;
        z-index: 1000;
        animation: slideIn 0.3s ease-out;
        box-shadow: 0 4px 15px rgba(0,0,0,0.3);
    `;
    
    document.body.appendChild(notification);
    
    setTimeout(() => {
        notification.style.animation = 'slideOut 0.3s ease-out';
        setTimeout(() => {
            if (notification.parentNode) {
                document.body.removeChild(notification);
            }
        }, 300);
    }, 3000);
}

// Check if user is logged in, redirect to login if not
async function ensureAuthenticated() {
    const currentPage = window.location.pathname.split('/').pop() || 'login.html';
    console.log('Current page:', currentPage);
    console.log('Auth token exists:', !!authToken);
    
    if (currentPage === 'login.html') {
        return;
    }
    
    if (!authToken) {
        console.log('No auth token, redirecting to login...');
        window.location.replace('login.html');
        throw new Error('Not authenticated');
    }
    
    try {
        await fetchAPI('/find-games', {}, true);
        console.log('Token validated successfully');
    } catch (error) {
        console.log('Token validation failed, redirecting to login...');
        localStorage.removeItem('authToken');
        localStorage.removeItem('username');
        authToken = null;
        window.location.replace('login.html');
        throw new Error('Token invalid');
    }
}

// Fetch games from API with optional filters
async function loadGames(search = '', publisher = '', genre = '') {
    try {
        await ensureAuthenticated();
        const params = new URLSearchParams();
        if (search) params.append('name', search);
        if (publisher) params.append('publisher', publisher);
        if (genre && genre !== 'All') params.append('genre', genre);
        
        const url = `/find-games${params.toString() ? '?' + params.toString() : ''}`;
        const data = await fetchAPI(url, {}, true);
        const games = data.games || [];
        renderGames(games);
    } catch (error) {
        console.error('Failed to load games:', error);
        renderGames([]);
    }
}

// Display games in the store grid
function renderGames(games) {
    const gamesGrid = document.getElementById('storeGames') || document.querySelector('.games-grid');
    if (!gamesGrid) return;
    
    if (games.length === 0) {
        gamesGrid.innerHTML = '<p class="empty-state">No games found</p>';
        return;
    }
    
    gamesGrid.innerHTML = games.map(game => `
        <div class="game-card" data-game-id="${game.ID}">
            <div class="game-image game-placeholder">
                ${game.name || 'Game'}
            </div>
            <div class="game-info">
                <h3 class="game-title">${game.name}</h3>
                <div class="game-tags">
                    <span class="tag">${game.genre || 'Game'}</span>
                    ${game.currency ? `<span class="tag">${game.currency}</span>` : ''}
                </div>
                <p class="game-price">${game.price ? formatPrice(game.price) : 'Free'}</p>
                ${document.getElementById('storeGames') ? `<button class="btn add-to-cart-btn" data-game-id="${game.ID}">Add to Cart</button>` : ''}
            </div>
        </div>
    `).join('');
    
    if (document.getElementById('storeGames')) {
        attachCartButtonListeners();
    }
}

function attachCartButtonListeners() {
    const buttons = document.querySelectorAll('.add-to-cart-btn');
    buttons.forEach(button => {
        button.addEventListener('click', async (e) => {
            e.stopPropagation();
            const gameId = parseInt(button.dataset.gameId);
            await addToCart(gameId);
        });
    });
}

// Add game to cart and refresh if on cart page
async function addToCart(gameId) {
    try {
        await ensureAuthenticated();
        await fetchAPI('/add-to-cart', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ game_id: gameId })
        }, true);
        showNotification('Game added to cart!');
        await updateCartCount();
        
        if (window.location.pathname.includes('cart.html')) {
            await loadCart();
        }
    } catch (error) {
        // Error already shown by fetchAPI
    }
}

async function updateCartCount() {
    try {
        await ensureAuthenticated();
        const data = await fetchAPI('/get-cart', {}, true);
        cart = data.cart_items || [];
    } catch (error) {
        console.error('Failed to update cart:', error);
    }
}

// Load user's purchased games and family shared games
async function loadLibrary() {
    try {
        await ensureAuthenticated();
        const data = await fetchAPI('/library', {}, true);
        const purchasedGames = data.games || [];
        
        renderLibrary(purchasedGames);
        await loadFamilyGames(purchasedGames);
    } catch (error) {
        console.error('Failed to load library:', error);
        renderLibrary([]);
    }
}

// Get games owned by family members (excluding user's own games)
async function loadFamilyGames(userGames = []) {
    try {
        const data = await fetchAPI('/get-family-games', {}, true);
        const allFamilyGames = data.games || [];
        const userGameIds = userGames.map(g => g.ID);
        const familyOnlyGames = allFamilyGames.filter(game => !userGameIds.includes(game.ID));
        
        renderFamilyGames(familyOnlyGames);
    } catch (error) {
        console.error('Failed to load family games:', error);
        renderFamilyGames([]);
    }
}

function renderFamilyGames(familyGames) {
    const familyGamesGrid = document.getElementById('familyGamesGrid');
    if (!familyGamesGrid) return;
    
    if (familyGames.length === 0) {
        familyGamesGrid.innerHTML = '<p class="empty-state empty-state-grid">No family games available. Family members\' games will appear here (excluding ones you own).</p>';
        return;
    }
    
    familyGamesGrid.innerHTML = familyGames.map(game => `
        <div class="library-item" data-game-id="${game.ID}">
            <div class="library-image library-placeholder">${game.name}</div>
            <div class="library-details">
                <h3 class="library-title">${game.name}</h3>
                <p class="library-playtime">${game.genre} • ${formatPrice(game.price)}</p>
                <p class="status-label status-family">Family Shared</p>
            </div>
        </div>
    `).join('');
}

function renderLibrary(libraryItems) {
    const libraryGrid = document.querySelector('.library-grid');
    if (!libraryGrid) return;
    
    if (libraryItems.length === 0) {
        libraryGrid.innerHTML = '<p class="empty-state empty-state-grid">Your library is empty. Visit the store to purchase games!</p>';
        return;
    }
    
    libraryGrid.innerHTML = libraryItems.map(game => `
        <div class="library-item" data-game-id="${game.ID}">
            <div class="library-image library-placeholder">${game.name}</div>
            <div class="library-details">
                <h3 class="library-title">${game.name}</h3>
                <p class="library-playtime">${game.genre} • ${formatPrice(game.price)}</p>
                <p class="status-label status-owned">Owned</p>
            </div>
        </div>
    `).join('');
    
    const playButtons = document.querySelectorAll('.play-btn');
    playButtons.forEach(button => {
        button.addEventListener('click', (e) => {
            const libraryItem = e.target.closest('.library-item');
            const gameTitle = libraryItem.querySelector('.library-title').textContent;
            showNotification(`Launching ${gameTitle}...`);
        });
    });
}

// Load user profile data and family members
async function loadProfile() {
    try {
        await ensureAuthenticated();
        
        const username = localStorage.getItem('username') || 'User';
        
        const profileUsername = document.getElementById('profileUsername');
        if (profileUsername) profileUsername.textContent = username;
        
        const avatarInitial = document.getElementById('avatarInitial');
        if (avatarInitial) avatarInitial.textContent = username.charAt(0).toUpperCase();
        
        const libraryData = await fetchAPI('/library', {}, true);
        const gamesCount = libraryData.games ? libraryData.games.length : 0;
        const gamesOwned = document.getElementById('gamesOwned');
        if (gamesOwned) gamesOwned.textContent = gamesCount;
        
        try {
            const familyData = await fetchAPI('/get-family', {}, true);
            const familySize = document.getElementById('familySize');
            const familyList = document.getElementById('familyList');
            const createFamilyBtn = document.getElementById('createFamilyBtn');
            const addFamilySection = document.getElementById('addFamilySection');
            
            if (familyData.family_members && familyData.family_members.length > 0) {
                const members = familyData.family_members;
                if (familySize) familySize.textContent = members.length;
                
                if (familyList) {
                    familyList.innerHTML = members.map(member => `
                        <div class="friend-card">
                            <div class="friend-avatar">${member.username.charAt(0).toUpperCase()}</div>
                            <div class="friend-info">
                                <div class="friend-name">${member.username}</div>
                                <div class="friend-status">${member.email}</div>
                            </div>
                        </div>
                    `).join('');
                }
                
                if (createFamilyBtn) createFamilyBtn.style.display = 'none';
                if (addFamilySection) addFamilySection.style.display = 'flex';
            } else {
                if (familySize) familySize.textContent = '0';
                if (familyList) {
                    familyList.innerHTML = '<p class="family-empty">Not in a family yet. Create one to get started!</p>';
                }
                
                if (createFamilyBtn) createFamilyBtn.style.display = 'block';
                if (addFamilySection) addFamilySection.style.display = 'none';
            }
        } catch (error) {
            console.log('Not in a family:', error);
            const familySize = document.getElementById('familySize');
            const createFamilyBtn = document.getElementById('createFamilyBtn');
            const addFamilySection = document.getElementById('addFamilySection');
            
            if (familySize) familySize.textContent = '0';
            if (createFamilyBtn) createFamilyBtn.style.display = 'block';
            if (addFamilySection) addFamilySection.style.display = 'none';
        }
        
        const createFamilyBtn = document.getElementById('createFamilyBtn');
        if (createFamilyBtn) {
            createFamilyBtn.addEventListener('click', createFamily);
        }
        
        const addToFamilyBtn = document.getElementById('addToFamilyBtn');
        if (addToFamilyBtn) {
            addToFamilyBtn.addEventListener('click', addUserToFamily);
        }
        
        const leaveFamilyBtn = document.getElementById('leaveFamilyBtn');
        if (leaveFamilyBtn) {
            leaveFamilyBtn.addEventListener('click', leaveFamily);
        }
        
    } catch (error) {
        console.error('Failed to load profile:', error);
    }
}

// Create a new family for the current user
async function createFamily() {
    try {
        await ensureAuthenticated();
        await fetchAPI('/create-family', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' }
        }, true);
        
        showNotification('Family created successfully!');
        await loadProfile();
    } catch (error) {
        // Error already shown by fetchAPI
    }
}

// Add another user to your family by username
async function addUserToFamily() {
    const usernameInput = document.getElementById('usernameToAdd');
    const username = usernameInput?.value.trim();
    
    if (!username) {
        showNotification('Please enter a username');
        return;
    }
    
    try {
        await ensureAuthenticated();
        const result = await fetchAPI('/add-user-to-family', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username: username })
        }, true);
        
        showNotification(result.result || 'User added to family!');
        usernameInput.value = '';
        await loadProfile();
    } catch (error) {}
}

// Leave current family with confirmation
async function leaveFamily() {
    if (!confirm('Are you sure you want to leave your family?')) {
        return;
    }
    
    try {
        await ensureAuthenticated();
        await fetchAPI('/leave-family', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' }
        }, true);
        
        showNotification('You have left the family');
        await loadProfile();
    } catch (error) {
        // Error already shown by fetchAPI
    }
}

// Load cart items and recommended games
async function loadCart() {
    try {
        await ensureAuthenticated();
        const data = await fetchAPI('/get-cart', {}, true);
        cart = data.cart_items || [];
        renderCart(cart);
        await loadRecommendedGames();
    } catch (error) {
        console.error('Failed to load cart:', error);
        renderCart([]);
    }
}

// Show 3 random games in cart recommendations
async function loadRecommendedGames() {
    try {
        const data = await fetchAPI('/find-games', {}, true);
        const allGames = data.games || [];
        const shuffled = allGames.sort(() => 0.5 - Math.random());
        const recommended = shuffled.slice(0, 3);
        
        renderRecommendedGames(recommended);
    } catch (error) {
        console.error('Failed to load recommended games:', error);
    }
}

function renderRecommendedGames(games) {
    const recommendedGrid = document.getElementById('recommendedGames');
    if (!recommendedGrid) return;
    
    if (games.length === 0) {
        recommendedGrid.innerHTML = '<p class="empty-state">No recommendations available</p>';
        return;
    }
    
    recommendedGrid.innerHTML = games.map(game => `
        <div class="game-card" data-game-id="${game.ID}">
            <div class="game-image game-placeholder">${game.name || 'Game'}</div>
            <div class="game-info">
                <h3 class="game-title">${game.name}</h3>
                <div class="game-tags">
                    <span class="tag">${game.genre || 'Game'}</span>
                    ${game.currency ? `<span class="tag">${game.currency}</span>` : ''}
                </div>
                <p class="game-price">${game.price ? formatPrice(game.price) : 'Free'}</p>
                <button class="btn add-to-cart-btn" data-game-id="${game.ID}">Add to Cart</button>
            </div>
        </div>
    `).join('');
    
    // Attach add to cart listeners
    const buttons = recommendedGrid.querySelectorAll('.add-to-cart-btn');
    buttons.forEach(button => {
        button.addEventListener('click', async (e) => {
            e.stopPropagation();
            const gameId = parseInt(button.dataset.gameId);
            await addToCart(gameId);
        });
    });
}

// Display cart items and calculate totals
function renderCart(cartItems) {
    const cartItemsContainer = document.querySelector('.cart-items');
    const cartSummaryContainer = document.querySelector('.cart-summary');
    
    if (!cartItemsContainer) return;
    
    if (cartItems.length === 0) {
        cartItemsContainer.innerHTML = '<p class="empty-state">Your cart is empty. Visit the store to add games!</p>';
        if (cartSummaryContainer) {
            cartSummaryContainer.style.display = 'none';
        }
        return;
    }
    
    cartItemsContainer.innerHTML = cartItems.map(game => `
        <div class="cart-item" data-game-id="${game.ID}">
            <div class="cart-item-image cart-placeholder">${game.name}</div>
            <div class="cart-item-details">
                <h3 class="cart-item-title">${game.name}</h3>
                <div class="game-tags">
                    <span class="tag">${game.genre}</span>
                    ${game.currency ? `<span class="tag">${game.currency}</span>` : ''}
                </div>
            </div>
            <div class="cart-item-price">${formatPrice(game.price)}</div>
            <button class="cart-item-remove" data-game-id="${game.ID}">×</button>
        </div>
    `).join('');
    
    let subtotal = cartItems.reduce((sum, game) => sum + (game.price || 0), 0);
    let discount = subtotal * 0.10;
    let tax = (subtotal - discount) * 0.08;
    let total = subtotal - discount + tax;
    
    if (cartSummaryContainer) {
        cartSummaryContainer.style.display = 'block';
        const summaryRows = cartSummaryContainer.querySelectorAll('.summary-value');
        if (summaryRows.length >= 4) {
            summaryRows[0].textContent = formatPrice(subtotal);
            summaryRows[1].textContent = formatPrice(tax);
            summaryRows[2].textContent = `-${formatPrice(discount)}`;
            summaryRows[3].textContent = formatPrice(total);
        }
    }
    
    const removeButtons = document.querySelectorAll('.cart-item-remove');
    removeButtons.forEach(button => {
        button.addEventListener('click', async (e) => {
            const gameId = button.dataset.gameId;
            await removeFromCart(gameId);
        });
    });
    
    const purchaseBtn = document.querySelector('.purchase-btn');
    if (purchaseBtn) {
        purchaseBtn.onclick = checkout;
    }
}

// Remove item from cart
async function removeFromCart(gameId) {
    try {
        await ensureAuthenticated();
        await fetchAPI('/remove-from-cart', { 
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ game_id: parseInt(gameId) })
        }, true);
        showNotification('Item removed from cart');
        await loadCart();
    } catch (error) {
        // Error already shown by fetchAPI
    }
}

// Process purchase and redirect to library
async function checkout() {
    try {
        await ensureAuthenticated();
        const result = await fetchAPI('/purchase', { 
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                street: '123 Main St',
                city: 'Gaming City',
                state: 'GC',
                zip_code: '12345',
                country: 'USA',
                payment_method: 'card'
            })
        }, true);
        showNotification('Purchase successful!');
        await loadCart();
        
        setTimeout(() => {
            window.location.href = 'library.html';
        }, 2000);
    } catch (error) {
        // Error already shown by fetchAPI
    }
}

// Filter games by search term
function setupSearch() {
    const searchInput = document.getElementById('searchInput');
    if (searchInput) {
        searchInput.addEventListener('input', function(e) {
            const searchTerm = e.target.value.toLowerCase();
            const gameCards = document.querySelectorAll('.game-card');
            
            gameCards.forEach(card => {
                const title = card.querySelector('.game-title').textContent.toLowerCase();
                card.style.display = title.includes(searchTerm) ? 'block' : 'none';
            });
        });
    }
}

// Category button click handlers for store
function setupCategoryFilters() {
    const categoryButtons = document.querySelectorAll('.category-btn');
    categoryButtons.forEach(button => {
        button.addEventListener('click', function() {
            categoryButtons.forEach(btn => btn.classList.remove('active'));
            this.classList.add('active');
            const genre = this.textContent === 'All' ? '' : this.textContent;
            loadGames('', '', genre);
        });
    });
}

// Filter buttons for library page
function setupLibraryFilters() {
    const filterButtons = document.querySelectorAll('.filter-btn');
    filterButtons.forEach(button => {
        button.addEventListener('click', function() {
            filterButtons.forEach(btn => btn.classList.remove('active'));
            this.classList.add('active');
            console.log('Filter selected:', this.textContent);
            // Could implement filtering logic here
        });
    });
}

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
    
    .library-playtime {
        color: #8f98a0;
        font-size: 0.9em;
        margin: 5px 0;
    }
`;
document.head.appendChild(style);

document.addEventListener('DOMContentLoaded', async function() {
    console.log('WaterVapor initializing...');
    
    const currentPage = window.location.pathname.split('/').pop() || 'login.html';
    
    if (currentPage === 'login.html') {
        console.log('Login page loaded');
        return;
    }
    
    await ensureAuthenticated();
    setupLoginUI();
    
    switch(currentPage) {
        case 'store.html':
            await loadGames();
            setupCategoryFilters();
            setupSearch();
            break;
        case 'library.html':
            await loadLibrary();
            setupLibraryFilters();
            break;
        case 'cart.html':
            await loadCart();
            break;
        case 'profile.html':
            await loadProfile();
            break;
        case 'landingpage.html':
        default:
            await loadGames();
            break;
    }
    
    document.querySelectorAll('a[href^="#"]').forEach(anchor => {
        anchor.addEventListener('click', function (e) {
            e.preventDefault();
            const target = document.querySelector(this.getAttribute('href'));
            if (target) {
                target.scrollIntoView({ behavior: 'smooth' });
            }
        });
    });
    
    console.log('WaterVapor initialized successfully!');
});

// Update header to show login/logout button
function setupLoginUI() {
    const userSections = document.querySelectorAll('.user-section');
    userSections.forEach(section => {
        const username = localStorage.getItem('username') || null;
        if (username) {
            section.innerHTML = `
                <div class="wallet">Welcome, ${username}!</div>
                <button class="btn btn-small" onclick="logout()">Logout</button>
            `;
        } else {
            section.innerHTML = `
                <button class="btn btn-small" onclick="showLoginModal()">Login</button>
            `;
        }
    });
}

// Display login popup modal
function showLoginModal() {
    const modal = document.createElement('div');
    modal.id = 'loginModal';
    modal.className = 'login-modal';
    
    modal.innerHTML = `
        <div class="login-modal-content">
            <h2>Login to WaterVapor</h2>
            <input type="text" id="loginUsername" placeholder="Username" />
            <input type="password" id="loginPassword" placeholder="Password" />
            <div class="login-modal-actions">
                <button onclick="performLogin()" class="btn">Login</button>
                <button onclick="closeLoginModal()" class="btn btn-cancel">Cancel</button>
            </div>
            <p class="login-modal-hint">Test users: alice, bob, charlie, diana, eve | Password: password123</p>
        </div>
    `;
    
    document.body.appendChild(modal);
    document.getElementById('loginUsername').focus();
    
    modal.addEventListener('keypress', (e) => {
        if (e.key === 'Enter') performLogin();
    });
}

function closeLoginModal() {
    const modal = document.getElementById('loginModal');
    if (modal) modal.remove();
}

// Handle login form submission
async function performLogin() {
    const username = document.getElementById('loginUsername').value;
    const password = document.getElementById('loginPassword').value;
    
    if (!username || !password) {
        showNotification('Please enter username and password', 'error');
        return;
    }
    
    try {
        const response = await fetch('/api/login', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({
                username: username,
                email: username + '@watervapor.com',
                password: password
            })
        });
        
        const data = await response.json();
        
        if (response.ok && data.token) {
            authToken = data.token;
            localStorage.setItem('authToken', authToken);
            localStorage.setItem('username', username);
            showNotification('Login successful!', 'success');
            closeLoginModal();
            setupLoginUI();
            location.reload();
        } else {
            showNotification(data.error || 'Login failed', 'error');
        }
    } catch (error) {
        showNotification('Login failed: ' + error.message, 'error');
    }
}

// Clear session and redirect to login
function logout() {
    localStorage.removeItem('authToken');
    localStorage.removeItem('username');
    authToken = null;
    window.location.href = 'login.html';
}