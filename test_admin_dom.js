const { JSDOM } = require('jsdom');
const fs = require('fs');

const html = fs.readFileSync('/Users/admin/workspace/jyq/wtb/admin-web/index.html', 'utf8');
const css = fs.readFileSync('/Users/admin/workspace/jyq/wtb/admin-web/app.css', 'utf8');

const dom = new JSDOM(html, {
  url: 'http://localhost:3000',
  runScripts: 'dangerously',
  resources: 'usable',
  pretendToBeVisual: true,
});

const window = dom.window;
const document = window.document;

window.localStorage = {
  _data: {},
  getItem(k) { return this._data[k] || null; },
  setItem(k, v) { this._data[k] = String(v); },
  removeItem(k) { delete this._data[k]; }
};

window.fetch = async (url, opts) => {
  if (url.includes('/admin/login')) {
    return { ok: true, status: 200, json: async () => ({ code: 200, data: { token: 'test_token' } }) };
  }
  return { ok: true, status: 200, json: async () => ({ code: 200, data: {} }) };
};

const style = document.createElement('style');
style.textContent = css;
document.head.appendChild(style);

const js = fs.readFileSync('/Users/admin/workspace/jyq/wtb/admin-web/app.js', 'utf8');
const script = document.createElement('script');
script.textContent = js;
document.body.appendChild(script);

setTimeout(async () => {
  console.log('=== Before login ===');
  const lp = document.getElementById('login-page');
  const app = document.getElementById('app');
  console.log('login-page:', lp.className, '| display:', window.getComputedStyle(lp).display);
  console.log('app:', app.className, '| display:', window.getComputedStyle(app).display);
  
  document.querySelectorAll('.modal').forEach(m => {
    console.log('modal', m.id, ':', m.className, '| display:', window.getComputedStyle(m).display);
  });
  
  await window.doLogin();
  await new Promise(r => setTimeout(r, 300));
  
  console.log('\n=== After login ===');
  console.log('login-page:', lp.className, '| display:', window.getComputedStyle(lp).display);
  console.log('app:', app.className, '| display:', window.getComputedStyle(app).display);
  
  document.querySelectorAll('.modal').forEach(m => {
    console.log('modal', m.id, ':', m.className, '| display:', window.getComputedStyle(m).display);
  });
  
  const db = document.getElementById('page-dashboard');
  console.log('dashboard:', db.className, '| display:', window.getComputedStyle(db).display);
  
  process.exit(0);
}, 100);
