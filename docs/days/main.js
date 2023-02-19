const monthName = [
  'Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun',
  'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];

const periodLength = 14;
const periodShift = 3

function main() {
  const today = dateToDay(new Date());
  const baseday = today - (today + periodShift) % periodLength;
  const table = document.querySelector('#body');
  let prevMonthName = '';
  for (let b = baseday - periodLength * 3; b < baseday + 350; b += periodLength) {
    let ymd;
    for (let d = 0; d < periodLength; d++) {
      const td = document.createElement('div');
      ymd = dayToDate(b + d);
      td.innerText = ymd.d;
      td.setAttribute('title', x(ymd));
      td.classList.add('cell');
      td.onclick = update;
      if (b + d === today) {
        td.classList.add('today');
      }
      if (d % 7 >= 5) {
        td.classList.add('dayoff');
      }
      table.appendChild(td);
    }
    let mn = monthName[ymd.m];
    if (mn === prevMonthName) {
      mn = '';
    } else {
      prevMonthName = mn;
    }
    const td = document.createElement('div'); // TODO function div builder
    td.innerText = mn;
    td.classList.add('cell');
    table.appendChild(td);
  }
}

function update() { // TODO
  console.log(this);
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

function x(d) {
  return `${(`${d.y}`).substring(2)}/${(`${101 + d.m}`).substring(1)}/${(`${100 + d.d}`).substring(1)}`;
}

main();
