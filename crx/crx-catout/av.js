// avito jira
(() => {
  const F = () => {
    const ee = document.getElementsByClassName('ghx-extra-field');
    [...ee].forEach(x => {
      const ttip = x.getAttribute('data-tooltip');
      if (!ttip) {
        return;
      }
      const tt = ttip.split(/:\s+/);
      if (tt[0] != 'Labels') {
        return;
      }
      const label = tt[1];
      if (!label) {
        return;
      }
      let p = 0;
      for (let i = 0; i < label.length; i++) {
        p *= 2;
        p += label.charCodeAt(i);
      }
      // const clr = '#' + (p % 4096).toString(16).padStart(3, '0');
      const clr = `hsl(${p % 360}, 100%, 30%)`;
      x.style.backgroundColor = clr;
      x.style.color = '#fff';
      x.style.fontSize = '10px';
      x.style.fontWeight = 'bold';
      x.style.borderRadius = '100%';
      x.style.padding = '2px 10px';
    });
    [...document.getElementsByTagName('aui-badge')].forEach(x => {
      let its = x.innerText;
      if (!its) {
        its = '7'; // fake
      }
      let it = +its;
      if (!it) { // NaN
        it = 7; // fake
      }
      // const clr = '#' + (it % 16) + '00';
      const clr = `hsl(${120 + (Math.floor((it / 8) * 240) % 240)}, 100%, 30%)`; // 120->360
      x.style.backgroundColor = clr;
      x.style.color = '#fff';
    });
  };
  F();
  setTimeout(F, 500);
  setTimeout(F, 1000);
  setInterval(F, 2000);
})();
