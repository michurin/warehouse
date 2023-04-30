// avito jira
(() => {
  const F = () => {
    const ee = document.getElementsByClassName('ghx-extra-field-content');
    [...ee].forEach(x => {
      if (x.childNodes && x.childNodes[0] && x.childNodes[0].nodeType === 1) {
        return; // child is not text (3): node already fixed
      }
      const itext = x.innerText;
      if (itext === 'None') {
        x.parentNode.removeChild(x);
        return
      }
      x.innerText = '';
      itext.split(/\s*,\s*/).forEach(label => {
        const s = document.createElement('span');
        s.innerText = label;
        let p = 0;
        for (let i = 0; i < label.length; i++) {
          p *= 2;
          p += label.charCodeAt(i);
        }
        // const clr = '#' + (p % 4096).toString(16).padStart(3, '0');
        const clr = 'hsl(' + (p % 360) + ', 100%, 30%)';
        s.style.backgroundColor = clr;
        s.style.color = '#fff';
        s.style.fontSize = '10px';
        s.style.fontWeight = 'bold';
        s.style.borderRadius = '100%';
        s.style.padding = '2px 10px';
        x.appendChild(s);
      });
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
  setTimeout(F, 50);
  setTimeout(F, 500);
  setInterval(F, 1000);
})();
