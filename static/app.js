let currentUserId = localStorage.getItem('userId');
let currentUsername = localStorage.getItem('username');
let jwtToken = localStorage.getItem('token');

// DOM Elements
const loginSection = document.getElementById('login-section');
const dashboardSection = document.getElementById('dashboard-section');
const loginForm = document.getElementById('login-form');
const registerForm = document.getElementById('register-form');
const addModal = document.getElementById('add-modal');
const vaultGrid = document.getElementById('vault-grid');
const notification = document.getElementById('notification');

// Navigation
document.getElementById('show-register').onclick = (e) => {
    e.preventDefault();
    loginForm.classList.add('hidden');
    registerForm.classList.remove('hidden');
};

document.getElementById('show-login').onclick = (e) => {
    e.preventDefault();
    registerForm.classList.add('hidden');
    loginForm.classList.remove('hidden');
};

document.getElementById('add-pwd-btn').onclick = () => addModal.classList.remove('hidden');
document.querySelector('.close-btn').onclick = () => addModal.classList.add('hidden');

// Auth logic
loginForm.onsubmit = async (e) => {
    e.preventDefault();
    const username = document.getElementById('login-username').value;
    const password = document.getElementById('login-password').value;

    try {
        const res = await fetch('/api/login', {
            method: 'POST',
            body: JSON.stringify({ username, password })
        });
        const data = await res.json();
        
        if (res.ok) {
            currentUserId = data.userId;
            currentUsername = username;
            jwtToken = data.token;
            
            localStorage.setItem('userId', currentUserId);
            localStorage.setItem('username', currentUsername);
            localStorage.setItem('token', jwtToken);

            showDashboard();
            notify("Login successful!");
        } else {
            notify(data || "Invalid credentials", true);
        }
    } catch (err) {
        notify("Server error", true);
    }
};

registerForm.onsubmit = async (e) => {
    e.preventDefault();
    const username = document.getElementById('reg-username').value;
    const password = document.getElementById('reg-password').value;

    try {
        const res = await fetch('/api/register', {
            method: 'POST',
            body: JSON.stringify({ username, password })
        });
        if (res.ok) {
            notify("Registration successful! You can now login.");
            document.getElementById('show-login').click();
        } else {
            const err = await res.text();
            notify(err, true);
        }
    } catch (err) {
        notify("Server error", true);
    }
};

document.getElementById('logout-btn').onclick = () => {
    currentUserId = null;
    currentUsername = "";
    jwtToken = null;
    localStorage.clear();
    dashboardSection.classList.add('hidden');
    loginSection.classList.remove('hidden');
};

// Vault logic
async function showDashboard() {
    loginSection.classList.add('hidden');
    dashboardSection.classList.remove('hidden');
    document.getElementById('display-username').innerText = currentUsername;
    loadVault();
}

async function loadVault() {
    try {
        const res = await fetch(`/api/passwords`, {
            headers: { 'Authorization': `Bearer ${jwtToken}` }
        });
        if (res.status === 401) return handleUnauthorized();
        const passwords = await res.json();
        renderVault(passwords);
    } catch (err) {
        notify("Failed to load vault", true);
    }
}

function renderVault(passwords) {
    vaultGrid.innerHTML = '';
    if (!passwords || passwords.length === 0) {
        vaultGrid.innerHTML = '<p style="text-align:center; width:100%; color:var(--text-muted)">Your vault is empty</p>';
        return;
    }

    passwords.forEach(p => {
        const card = document.createElement('div');
        card.className = 'pwd-card';
        card.innerHTML = `
            <div class="card-header">
                <span class="card-label">${p.label}</span>
                <div class="action-btns">
                    <button class="icon-btn" onclick="toggleVisibility('${p.id}')">VIEW</button>
                    <button class="icon-btn" onclick="copyToClipboard('${p.password}')">COPY</button>
                </div>
            </div>
            <div class="pwd-display">
                <code id="pwd-${p.id}" data-real="${p.password}">••••••••</code>
            </div>
            <div class="card-footer">
                Added: ${new Date(p.created_at).toLocaleDateString()}
            </div>
        `;
        vaultGrid.appendChild(card);
    });
}

window.toggleVisibility = (id) => {
    const el = document.getElementById(`pwd-${id}`);
    if (el.innerText === '••••••••') {
        el.innerText = el.getAttribute('data-real');
    } else {
        el.innerText = '••••••••';
    }
};

// Add Password
let currentCriteria = null;

const pwdLabelInput = document.getElementById('pwd-label');
const pwdValueInput = document.getElementById('pwd-value');
const criteriaDisplay = document.getElementById('criteria-display');
const criteriaList = document.getElementById('criteria-list');
const pwdError = document.getElementById('pwd-error');

pwdLabelInput.oninput = async () => {
    const label = pwdLabelInput.value.trim();
    if (label.length < 2) {
        criteriaDisplay.classList.add('hidden');
        currentCriteria = null;
        return;
    }

    try {
        const res = await fetch(`/api/criteria?label=${encodeURIComponent(label)}`, {
            headers: { 'Authorization': `Bearer ${jwtToken}` }
        });
        if (res.status === 401) return handleUnauthorized();
        if (res.ok) {
            currentCriteria = await res.json();
            renderCriteria(currentCriteria);
            criteriaDisplay.classList.remove('hidden');
            validateWithUI();
        } else {
            criteriaDisplay.classList.add('hidden');
            currentCriteria = null;
        }
    } catch (err) {
        console.error("Failed to fetch criteria", err);
    }
};

pwdValueInput.oninput = () => {
    validateWithUI();
};

function renderCriteria(c) {
    criteriaList.innerHTML = '';
    const items = [
        { key: 'length', text: `Min Length: ${c.min_length}` },
        { key: 'upper', text: `Min Uppercase: ${c.min_uppercase}` },
        { key: 'lower', text: `Min Lowercase: ${c.min_lowercase}` },
        { key: 'number', text: `Min Numbers: ${c.min_numbers}` },
        { key: 'special', text: `Min Special: ${c.min_special} (${c.allowed_special})` }
    ];

    items.forEach(item => {
        const li = document.createElement('li');
        li.className = 'criteria-item';
        li.id = `crit-${item.key}`;
        li.innerHTML = `<span class="criteria-icon">○</span> ${item.text}`;
        criteriaList.appendChild(li);
    });
}

function validateWithUI() {
    if (!currentCriteria) return true;

    const pass = pwdValueInput.value;
    const c = currentCriteria;
    
    let upper = 0, lower = 0, number = 0, special = 0;
    for (const char of pass) {
        if (/[A-Z]/.test(char)) upper++;
        else if (/[a-z]/.test(char)) lower++;
        else if (/[0-9]/.test(char)) number++;
        else if (c.allowed_special.includes(char)) special++;
    }

    const checks = {
        length: pass.length >= c.min_length,
        upper: upper >= c.min_uppercase,
        lower: lower >= c.min_lowercase,
        number: number >= c.min_numbers,
        special: special >= c.min_special
    };

    let allValid = true;
    for (const key in checks) {
        const el = document.getElementById(`crit-${key}`);
        if (checks[key]) {
            el.classList.add('valid');
            el.classList.remove('invalid');
            el.querySelector('.criteria-icon').innerText = '✓';
        } else {
            el.classList.remove('valid');
            el.classList.add('invalid');
            el.querySelector('.criteria-icon').innerText = '○';
            allValid = false;
        }
    }

    if (pass.length > 0 && c.allowed_special.includes(pass[0])) {
        pwdError.innerText = "Password cannot start with a special character";
        pwdError.style.display = 'block';
        allValid = false;
    } else {
        pwdError.style.display = 'none';
    }

    return allValid;
}

document.getElementById('gen-pwd-btn').onclick = async () => {
    const label = document.getElementById('pwd-label').value;
    const res = await fetch(`/api/generate?label=${encodeURIComponent(label)}`, {
        headers: { 'Authorization': `Bearer ${jwtToken}` }
    });
    if (res.status === 401) return handleUnauthorized();
    const data = await res.json();
    document.getElementById('pwd-value').value = data.password;
    validateWithUI();
};

document.getElementById('add-pwd-form').onsubmit = async (e) => {
    e.preventDefault();
    
    if (currentCriteria && !validateWithUI()) {
        notify("Password does not meet requirements", true);
        return;
    }

    const label = document.getElementById('pwd-label').value;
    const password = document.getElementById('pwd-value').value;

    try {
        const res = await fetch('/api/passwords', {
            method: 'POST',
            headers: { 
                'Authorization': `Bearer ${jwtToken}`,
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ label, password })
        });
        if (res.status === 401) return handleUnauthorized();
        if (res.ok) {
            addModal.classList.add('hidden');
            document.getElementById('add-pwd-form').reset();
            criteriaDisplay.classList.add('hidden');
            currentCriteria = null;
            loadVault();
            notify("Password saved!");
        } else {
            const err = await res.text();
            notify(err, true);
        }
    } catch (err) {
        notify("Failed to save password", true);
    }
};

// Utils
function notify(msg, isError = false) {
    notification.innerText = msg;
    notification.style.background = isError ? 'var(--accent)' : 'var(--primary)';
    notification.classList.remove('hidden');
    setTimeout(() => notification.classList.add('hidden'), 3000);
}

window.copyToClipboard = (text) => {
    navigator.clipboard.writeText(text).then(() => notify("Copied to clipboard!"));
};

function handleUnauthorized() {
    notify("Session expired. Please login again.", true);
    document.getElementById('logout-btn').click();
}

// Auto-login if token exists
if (jwtToken && currentUserId) {
    showDashboard();
}
