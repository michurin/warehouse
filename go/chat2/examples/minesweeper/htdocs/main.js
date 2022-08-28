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

// ***** KIT *****

const kit = chatAdapter();

kit.loop((data) => {
  console.log('LOOP', data);
  if (data.chat) {
    data.chat.forEach((msg) => { // TODO check type is array
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
      console.log('GAME MSG', msg);
      if (msg.r) {
        initGameArena(msg.r.w, msg.r.h);
        userTable = {};
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
        setTimeout(() => alert('Game Over'), 10);
      }
    });
    const tp = $('#top');
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

$(() => {
  const text = $('#text');
  const name = $('#name');
  const color = $('#color');
  text.keypress((e) => {
    if (e.which === 13) {
      kit.send({
        color: color.val(),
        text: text.val(),
        name: name.val(),
      });
      text.val('');
      text.focus();
      return false;
    }
    return true;
  });
  name.keypress((e) => {
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

function buildOnClick(i, j, e) {
  return async () => {
    // $('#game').css('cursor', 'wait');
    console.log(i, j, e);
    kit.game({
      x: i,
      y: j,
      cid: CID,
      name: $('#name').val(),
      color: $('#color').val(),
    });
    // $('#game').css('cursor', 'default');
    $('#text').focus();
    return false;
  };
}

function buildOnRightClick(i, j, e) {
  return () => {
    console.log(i, j, e, e.text());
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
      tr.append(e.attr('id', id).click(buildOnClick(i, j, e)).contextmenu(buildOnRightClick(i, j, e)));
    }
    tbl.append(tr);
  }
  $('#game').empty().append(tbl);
}

// ***** SAVE/RESTORE *****

function setMessageColor(clr) {
  $('#name').css({ color: clr });
  $('#text').css({ color: clr });
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
    console.log('random color', v);
  }
  clr.val(v);
  setMessageColor(v);
  const name = $('#name');
  name.on('input', () => {
    // TODO validate!
    localStorage.setItem('name', name.val());
  });
  v = localStorage.getItem('name') || 'noname';
  if (v) {
    name.val(v); // TODO validate?
  }
});