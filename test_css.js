const { JSDOM } = require('jsdom');

const html = `
<!DOCTYPE html>
<html><head><style>
.hidden { display: none !important; }
.modal { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.4); display: flex; z-index: 100; }
.login-page { display: flex; min-height: 100vh; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); }
.app { display: flex; min-height: 100vh; }
</style></head><body>
<div id="login-page" class="login-page">LOGIN</div>
<div id="app" class="app hidden">APP</div>
<div id="modal1" class="modal hidden">MODAL</div>
</body></html>`;

const dom = new JSDOM(html, { pretendToBeVisual: true });
const window = dom.window;
const document = window.document;

console.log('=== Initial state ===');
console.log('login-page display:', window.getComputedStyle(document.getElementById('login-page')).display);
console.log('app display:', window.getComputedStyle(document.getElementById('app')).display);
console.log('modal1 display:', window.getComputedStyle(document.getElementById('modal1')).display);

// Simulate login
document.getElementById('login-page').classList.add('hidden');
document.getElementById('app').classList.remove('hidden');

console.log('\n=== After login ===');
console.log('login-page display:', window.getComputedStyle(document.getElementById('login-page')).display);
console.log('app display:', window.getComputedStyle(document.getElementById('app')).display);
console.log('modal1 display:', window.getComputedStyle(document.getElementById('modal1')).display);
