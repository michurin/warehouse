const bb = document.getElementsByTagName('button');
const ct = document.getElementsByTagName('p')[1];
const beep = new Audio('beep.mp3');

let tid;
let sv = 30;

function stop() {
  if (tid === undefined) {
    return;
  }
  clearInterval(tid);
  tid = undefined;
}

function step() {
  const v = (Number(ct.innerText) || 0) - 1;
  ct.innerText = v;
  if (v <= 0) {
    stop();
    beep.play();
    document.body.style.backgroundColor = '#fff'
    setTimeout(() => document.body.style.backgroundColor = '#000', 200);
  }
}

bb[0].onclick = () => { stop(); beep.play(); };
bb[1].onclick = () => { stop(); ct.innerText = sv = 30; };
bb[2].onclick = () => { stop(); ct.innerText = sv = 60; };
bb[3].onclick = () => { stop(); ct.innerText = sv; tid = setInterval(step, 1000); };
