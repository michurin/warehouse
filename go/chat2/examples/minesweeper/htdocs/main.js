// ***** CID *****

let CID = '';
let userTable = {};

$(() => {
  try {
    CID = localStorage.getItem('cid');
  } catch (e) { }
  CID = CID || '';
  if (CID.length !== 24) {
    CID = `${Date.now().toString(36)}|`;
    while (CID.length < 24) {
      CID += String.fromCharCode(Math.floor(Math.random() * 94) + 33);
    }
  }
  try {
    localStorage.setItem('cid', CID);
  } catch (e) { }
});

// ***** ALERT *****

$(() => {
  const note = $('#note');
  note.hide();
  note.click(() => note.hide());
});

var notificatoinTimeoutID;

function displayNotification(text, timeout) {
  const note = $('#note');
  note.text(text);
  note.show();
  clearTimeout(notificatoinTimeoutID);
  if (timeout) {
    notificatoinTimeoutID = setTimeout(hideNotification, timeout);
  }
}

function hideNotification() {
  $('#note').hide();
}

// ***** KIT *****

function checkRequest(o) { // very naive checker
  const k = Object.keys(o);
  for (let i = 0; i < k.length; i++) {
    const v = o[k[i]];
    if (typeof v === 'string' && v.length === 0) {
      return undefined;
    }
  }
  return o;
}

const kit = chatAdapter();

kit.loop((data) => {
  if (data.chat) {
    data.chat.forEach((msg) => {
      $('#board').append($('<div>').append(
        $('<b>').text(`${msg.name}:`),
        $('<span>').text(` ${msg.text}`),
      ).css({
        color: msg.color,
        'overflow-wrap': 'anywhere',
      }));
    });
    const c = $('#jail');
    c.scrollTop(c.prop('scrollHeight'));
    setTimeout(() => c.scrollTop(c.prop('scrollHeight')), 20);
  }
  if (data.game) {
    data.game.forEach((msg) => {
      console.log('GAME', msg); // TODO REMOVE
      if (msg.r) {
        displayNotification('Game started!', 1500);
        initGameArena(msg.r.w, msg.r.h);
        userTable = {}; // TODO save high scores for history
      }
      if (msg.f) {
        userTable = {};
      }
      if (msg.u) {
        msg.u.forEach((e) => {
          userTable[e.id] = e;
        });
      }
      // userTable must be updated before call setCell
      if (msg.f) {
        initGameArena(msg.f[0].length, msg.f.length);
        for (let j = 0; j < msg.f.length; j++) {
          const p = msg.f[j];
          for (let i = 0; i < p.length; i++) {
            setCell(i, j, p[i]);
          }
        }
      }
      if (msg.a) {
        msg.a.forEach((e) => { setCell(e.x, e.y, e.v); });
      }
      if (msg.go) {
        displayNotification('Game Over!\n\nYou can click on game area to start new one. Or just wait unitil someone click.');
      }
    });
    const tp = $('#top-content');
    tp.empty();
    const uu = Object.values(userTable);
    uu.sort((a, b) => b.score - a.score);
    uu.forEach((u) => {
      const st = { color: u.color };
      tp.append(
        $('<div>').text(u.name).css(st),
        $('<div>').text(u.score).css(st),
      );
    });
  }
});

// ***** CHAT *****

const nameLenLimit = 20;
const messageLenLimit = 200;

function filnalName() {
  const n = $('#name');
  const v = n.val()
    .substring(0, nameLenLimit)
    .replace(/[^\p{L}\p{M}\p{N}\p{P}\p{S}\p{Zs}]+/ug, '')
    .replace(/[^\p{L}\p{M}\p{N}\p{P}\p{S}]+/ug, ' ')
    .replace(/^[^\p{L}\p{M}\p{N}\p{P}\p{S}]+/ug, '')
    .replace(/[^\p{L}\p{M}\p{N}\p{P}\p{S}]+$/ug, '');
  if (v.length === 0) {
    return 'Incognito';
  }
  return v;
}

function fixName() {
  setTimeout(() => {
    const n = $('#name');
    n.val(n.val()
      .substring(0, nameLenLimit)
      .replace(/[^\p{L}\p{M}\p{N}\p{P}\p{S}\p{Zs}]+/ug, '')
      .replace(/[^\p{L}\p{M}\p{N}\p{P}\p{S}]+/ug, ' '));
    localStorage.setItem('name', filnalName());
  }, 200);
}

function getName() {
  const t = filnalName();
  $('#name').val(t);
  return t;
}

function fixMessage() {
  setTimeout(() => {
    const n = $('#text');
    n.val(n.val().substring(0, messageLenLimit).replace(/[^\p{L}\p{M}\p{N}\p{P}\p{S}\p{Zs}]+/ug, ''));
  }, 200);
}

function getMessage() {
  const n = $('#text');
  const t = n.val()
    .substring(0, messageLenLimit)
    .replace(/[^\p{L}\p{M}\p{N}\p{P}\p{S}\p{Zs}]+/ug, '')
    .replace(/[^\p{L}\p{M}\p{N}\p{P}\p{S}]+/ug, ' ')
    .replace(/^[^\p{L}\p{M}\p{N}\p{P}\p{S}]+/ug, '')
    .replace(/[^\p{L}\p{M}\p{N}\p{P}\p{S}]+$/ug, '');
  n.val(t);
  return t;
}

function getColor() {
  const clr = $('#color').val();
  if (/^#[0-9a-fA-F]{6}$/.test(clr)) {
    return clr;
  }
  return '';
}

$(() => {
  const text = $('#text');
  text.keypress((e) => {
    fixMessage();
    if (e.which === 13) {
      kit.send(checkRequest({
        color: getColor(),
        text: getMessage(),
        name: getName(),
      }));
      text.val('');
      text.focus();
      return false;
    }
    return true;
  });
  $('#name').keypress((e) => {
    fixName();
    if (e.which === 13) {
      text.focus();
      return false;
    }
    return true;
  });
  text.focus();
});

// ***** GAME *****

function setCell(x, y, e) {
  const uid = Math.floor(e / 10);
  const v = e % 10;
  const u = userTable[uid]; // TODO what if not found?
  if (!u) {
    return; // it's a bug on backend?
  }
  $(`#c${x}x${y}`).text(v === 9 ? 'W' : v).css({
    color: u.color,
    'background-color': v === 9 ? '#800' : '#000', // TODO contrast!
  }).prop('title', u.name);
}

function buildOnClick(i, j) {
  return async () => {
    // $('#game').css('cursor', 'wait');
    kit.game(checkRequest({
      x: i,
      y: j,
      cid: CID,
      name: getName(),
      color: getColor(),
    }));
    // $('#game').css('cursor', 'default');
    $('#text').focus();
    return false;
  };
}

function buildOnRightClick(e) {
  return () => {
    if (e.text() === '') {
      e.text('F').css('color', '#f00');
    } else if (e.text() === 'F') {
      e.text('');
    }
    $('#text').focus();
    return false;
  };
}

function initGameArena(w, h) {
  const tbl = $('<table>');
  for (let j = 0; j < h; j++) {
    const tr = $('<tr>');
    for (let i = 0; i < w; i++) {
      const id = `c${i}x${j}`;
      const e = $('<td>');
      tr.append(e.attr('id', id).click(buildOnClick(i, j)).contextmenu(buildOnRightClick(e)));
    }
    tbl.append(tr);
  }
  $('#game').empty().append(tbl);
}

// ***** SAVE/RESTORE *****

function setMessageColor(clr) {
  const t = clr.match(/^#([0-9a-fA-F]{2})([0-9a-fA-F]{2})([0-9a-fA-F]{2})$/);
  let r = .5 + .5 * Math.random();
  let g = .5 + .5 * Math.random();
  let b = .5 + .5 * Math.random();
  let k;
  if (t) {
    r = parseInt(t[1], 16);
    g = parseInt(t[2], 16);
    b = parseInt(t[3], 16);
  }
  const lum = 0.2126 * r + 0.7152 * g + 0.0722 * b;
  if (lum < 100) {
    if (lum < 1) {
      k = 100; // let lum=1
    } else {
      k = 100 / lum;
    }
    r *= k; // TODO naive
    g *= k;
    b *= k;
    clr = '#' + ([r, g, b].map(x => (Math.ceil(x > 255 ? 255 : x)).toString(16).padStart(2, '0')).join(''));
    $('#color').val(clr);
  }
  $('#name').css({ color: clr });
  $('#text').css({ color: clr });
  $('#color').parent().css('background-color', clr);
  localStorage.setItem('color', clr);
}

$(() => {
  let v;
  const clr = $('#color');
  clr.on('change', () => {
    setMessageColor(clr.val());
  });
  v = localStorage.getItem('color') || '';
  if (!/^#[0-9a-fA-F]{6}$/.test(v)) {
    v = '#';
    for (let i = 0; i < 3; i++) {
      v += Math.floor(320 + Math.random() * 192).toString(16).substring(1);
    }
  }
  clr.val(v);
  setMessageColor(v);
  v = fixName(localStorage.getItem('name')) || 'noname';
  $('#name').val(v);
});

// ***** MISC *****

$(window).bind('beforeunload', function() {
  return 'Are you sure you want to leave?';
});
