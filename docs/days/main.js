// generated

const classChain = {
  'dayoff': '',
  '': 'holiday',
  'holiday': 'vacation',
  'vacation': 'special',
  'special': 'dayoff',
  'dayoff-blue': 'blue',
  'blue': 'holiday-blue',
  'holiday-blue': 'vacation-blue',
  'vacation-blue': 'special-blue',
  'special-blue': 'dayoff-blue',
  'dayoff-gray': 'gray',
  'gray': 'holiday-gray',
  'holiday-gray': 'vacation-gray',
  'vacation-gray': 'special-gray',
  'special-gray': 'dayoff-gray',
  'dayoff-green': 'green',
  'green': 'holiday-green',
  'holiday-green': 'vacation-green',
  'vacation-green': 'special-green',
  'special-green': 'dayoff-green',
  'dayoff-red': 'red',
  'red': 'holiday-red',
  'holiday-red': 'vacation-red',
  'vacation-red': 'special-red',
  'special-red': 'dayoff-red',
  'dayoff-yellow': 'yellow',
  'yellow': 'holiday-yellow',
  'holiday-yellow': 'vacation-yellow',
  'vacation-yellow': 'special-yellow',
  'special-yellow': 'dayoff-yellow',
};

const bgChain = {
  'yellow': '',
  '': 'blue',
  'blue': 'gray',
  'gray': 'green',
  'green': 'red',
  'red': 'yellow',
  'holiday-yellow': 'holiday',
  'holiday': 'holiday-blue',
  'holiday-blue': 'holiday-gray',
  'holiday-gray': 'holiday-green',
  'holiday-green': 'holiday-red',
  'holiday-red': 'holiday-yellow',
  'vacation-yellow': 'vacation',
  'vacation': 'vacation-blue',
  'vacation-blue': 'vacation-gray',
  'vacation-gray': 'vacation-green',
  'vacation-green': 'vacation-red',
  'vacation-red': 'vacation-yellow',
  'special-yellow': 'special',
  'special': 'special-blue',
  'special-blue': 'special-gray',
  'special-gray': 'special-green',
  'special-green': 'special-red',
  'special-red': 'special-yellow',
  'dayoff-yellow': 'dayoff',
  'dayoff': 'dayoff-blue',
  'dayoff-blue': 'dayoff-gray',
  'dayoff-gray': 'dayoff-green',
  'dayoff-green': 'dayoff-red',
  'dayoff-red': 'dayoff-yellow',
};

// /generated

const monthName = [
  'Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun',
  'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];

const periodLength = 14;
const periodShift = 3;

const days = JSON.parse(window.localStorage.getItem('days') || '{}');

function main() {
  const today = dateToDay(new Date());
  const baseday = today - (today + periodShift) % periodLength;
  const table = document.querySelector('#body');
  let prevMonthName = '';
  for (let b = baseday - periodLength * 3; b < baseday + 350; b += periodLength) {
    let ymd;
    for (let d = 0; d < periodLength; d++) {
      const s = b + d;
      const td = document.createElement('div');
      ymd = dayToDate(s);
      td.innerText = ymd.d;
      td.setAttribute('title', `${(`${ymd.y}`).substring(2)}/${(`${101 + ymd.m}`).substring(1)}/${(`${100 + ymd.d}`).substring(1)}`);
      td.onclick = updater(s);
      if (s === today) {
        td.style.margin = 0;
        td.style.border = '2px dashed #000';
        td.style.animation = 'blinker 1s linear infinite';
        td.style.margin = '-5px';
        td.style.padding = '3px';
      }
      td.className = days[s] || (d % 7 >= 5 ? 'dayoff' : '');
      table.appendChild(td);
      const label = document.createElement('div');
      td.appendChild(label);
    }
    let mn = monthName[ymd.m];
    if (mn === prevMonthName) {
      mn = '';
    } else {
      prevMonthName = mn;
    }
    const td = document.createElement('div'); // TODO function div builder
    td.innerText = mn;
    table.appendChild(td);
  }
}

function updater(s) {
  return function (e) {
    e.preventDefault();
    if (document.getElementById('lock').checked) {
      return;
    }
    let chain = {}; // reset chain
    if (document.getElementById('bg').checked) {
      chain = bgChain;
    }
    if (document.getElementById('kind').checked) {
      chain = classChain;
    }
    const c = chain[this.className] || '';
    if (c === '') {
      delete (days[s]);
    } else {
      days[s] = c;
    }
    this.className = c;
    window.localStorage.setItem('days', JSON.stringify(days));
  };
}

function dateToDay(d) {
  return Date.UTC(d.getFullYear(), d.getMonth(), d.getDate()) / 86400000;
}

function dayToDate(n) {
  const t = new Date(n * 86400000);
  return {
    y: t.getFullYear(),
    m: t.getMonth(),
    d: t.getDate(),
  };
}

// autolock (to move to separate file?)

let tid;
function lock() {
  if (tid) {
    clearTimeout(tid);
  }
  tid = setTimeout(() => document.getElementById('lock').checked = true, 60000);
}
document.querySelectorAll('input[type=radio]').forEach((x) => x.onchange = lock);

// run

main();
