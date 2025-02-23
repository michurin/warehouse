// avito link in ms exchange
(() => {
  const re = /[^/]*:\/\//;
  const d2 = (x) => {
    if (x >= 10) {
      return x
    }
    return '0' + x;
  }
  let n = 0
  const F = () => {
    Array.from(document.querySelectorAll('*'))
      .map(x => (x.childNodes.length === 1 ? (x.childNodes[0].nodeType === 3 ? (x.innerText.match('^(http|avito.zoom.us/)') ? x : false) : false) : false))
      .filter(x => x)
      .forEach(x => {
        const c = x.childNodes[0];
        let t = x.innerText;
        if (!t.startsWith('http')) { // heuristic for naked domains
          t = `https://${t}`;
        }
        x.removeChild(c);
        const a = document.createElement('a');
        a.href = t;
        a.innerText = t.replace(re, '');
        a.target = '_blank';
        a.style.color = '#080';
        a.style.fontWeight = 'bold';
        x.append(a);
      });

    Array.from(document.getElementsByClassName('_wx_m1')).forEach(x => {
      x.style.borderStyle = 'solid';
      x.style.borderWidth = '1px 0 0 0';
    });

    Array.from(document.getElementsByClassName('_cb_G1 _cb_H1')).forEach(x => {
      x.style.background = 'linear-gradient(rgba(192, 224, 255, 1), rgba(255, 255, 255, 1))';
      const c = x.children;
      for (let i = 0; i < c.length; i++) { // recursion for poor
        c[i].style.background = 'none';
        const d = c[i].children;
        for (let j = 0; j < d.length; j++) {
          d[j].style.background = 'none';
        }
      }
    });

    let x = document.getElementsByClassName('_cb_u2 _cb_s2 ms-border-color-themeSecondary');
    if (x.length > 0) {
      x = x[0]
      x.style.borderWidth = '1px';
      x.style.borderStyle = 'solid';
      x.style.borderColor = 'rgba(255, 0, 0, .5)';
      x.style.zIndex = 100;
      let c = x.children;
      if (c.length > 0) {
        c = c[0]
      } else {
        c = document.createElement('div')
        c.style.backgroundColor = 'rgba(255, 0, 0, .2)';
        c.style.position = 'absolute';
        c.style.bottom = 0;
        c.style.right = 0;
        c.style.fontFamily = 'monospace'; //'sans-serif'
        c.style.fontSize = '9px';
        c.style.fontWeight = 'bold';
        c.style.padding = '0 .4em';
        c.style.borderRadius = '.5em .5em 0 0';
        c.style.color = '#330000';
        x.append(c)
      }
      const m = new Date;
      n = (n + 1) % 2;
      c.innerText = d2(m.getHours()) + [':', ' '][n] + d2(m.getMinutes());
    }
  };
  setTimeout(F, 500);
  setTimeout(F, 1000);
  setInterval(F, 2000);
})();
