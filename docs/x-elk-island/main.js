/*global document:readable, setTimeout:readable, localStorage:readable*/
/*eslint-env es6*/
/*eslint indent: ["error", 2]*/
/*eslint eqeqeq: ["error", "always"]*/
/*eslint prefer-const: "error"*/
/*eslint no-var: "error"*/
/*eslint no-undef: "error"*/
/*eslint one-var: ["error", "never"]*/
/*eslint semi: ["error", "never"]*/
/*eslint quotes: ["error", "single"]*/
/*eslint prefer-arrow-callback: "error"*/
/*eslint arrow-body-style: "error"*/

const debug = function() {
  let d = false
  window.addEventListener('keydown', (e) => { // eslint-disable-line no-undef
    if (e.key === '\x20' && e.target === document.body) {
      e.preventDefault()
    }
  })
  document.body.onkeyup = function(e) {
    if (e.key === '\x20') {
      d = !d
      document.getElementById('debug').style.display = d ? 'block' : 'none'
      render()
    }
  }
  return (v) => {
    if (v !== undefined) {
      d = v
    }
    return d
  }
}()

function log() {
  const t = []
  const s = []
  for (let i = 0; i < arguments.length; i++) {
    t.push(arguments[i])
    if (arguments[i].stack) {
      s.push(arguments[i].stack)
    }
  }
  const h = document.getElementById('log')
  const e = document.createElement('div')
  e.innerText = t.join(' ')
  h.prepend(e)
  s.forEach((x) => {
    const e = document.createElement('pre')
    e.innerText = x
    h.prepend(e)
  })
}

const arena = [] // it has to be split: data (mutable), pointers to DOM (immutable)
const gameState = {
  lines: 0,
  opens: 0,
  booms: 0,
}

const [lock, unlock, locked] = (() => {
  let l = false
  return [
    () => {
      l = true
      document.getElementById('game').style.cursor = 'wait'
    },
    () => {
      l = false
      document.getElementById('game').style.cursor = 'pointer'
    },
    () => l,
  ]
})()

function handler(o, rightClick) { // we must to use o instead (x, y) to manage shifted lines
  return (ev) => {
    ev.stopPropagation()
    ev.preventDefault()
    if (locked()) {
      log('LOCKED')
      return
    }
    const x = o.i
    const y = o.j
    log('click', x, y, rightClick)
    if (rightClick) {
      arena[y][x].flag = !arena[y][x].flag // .flag is not correlate to .open
    } else {
      // TODO check boom
      open(x, y)
      collapse() // ugly: collapse only if opened; return if no more collapsing
    }
    render() // ugly: update entire arena
    return false
  }
}

function collapse() {
  lock()
  for (let j = arena.length - 1; j >= 0; j--) {
    let f = true
    arena[j].forEach((q) => {
      f = f && (q.mine || q.open)
    })
    if (f) {
      delayed(1, evaporate, j) // we must delay it to avoid double rendering
      return // without unlocking
    }
  }
  unlock() // unlock if nothing more to do
}

function evaporate(k) {
  for (let i = 0; i < arena[k].length; i++) {
    const q = arena[k][i]
    q.flag = q.mine && !q.open // just show all mines
  }
  render() // double rendering: one in handler, and one here
  for (let i = 0; i < arena[k].length; i++) {
    const q = arena[k][i]
    if (!q.flag && !q.mine) {
      q.element.div.style.backgroundColor = '#fff'
    }
  }
  delayed(1000, evaporate_remove, k, 0)
}

function evaporate_remove(k, step) {
  if (step === 0) {
    for (let i = 0; i < arena[k].length; i++) {
      arena[k][i].element.div.remove()
    }
  }
  if (k === 0) { // hack: skip movement if nothing to move
    step = 10
  }
  for (let j = 0; j < k; j++) {
    const v = j + step / 10
    for (let i = 0; i < arena[j].length; i++) {
      arena[j][i].element.div.style.top = `calc(${v} * 4vmin)` // spaces are important
    }
  }
  if (step < 10) { // continue movement
    delayed(20, evaporate_remove, k, step + 1)
    return
  }
  // finishing movement
  for (let j = 0; j < k; j++) {
    for (let i = 0; i < arena[j].length; i++) {
      arena[j][i].j++
    }
  }
  gameState.lines++
  arena.splice(k, 1)
  arena.unshift(initLine())
  render()
  delayed(500, evaporate_open)
}

function evaporate_open() {
  for (let j = 0; j < arena.length; j++) { // touching every open cell to force zero-neighbors-propagation
    for (let i = 0; i < arena[j].length; i++) {
      if (arena[j][i].open) {
        open(i, j) // can be optimized
      }
    }
  }
  render()
  delayed(500, collapse) // next collapse; TODO delay only if something was opened
}

function open(x, y) {
  if (!arena[y][x].open) { // slightly hackish extra check
    gameState.opens++
    if (arena[y][x].mine) {
      gameState.booms++
    }
  }
  arena[y][x].open = true
  if (arena[y][x].mine) {
    return // TODO boom
  }
  const n = neighbours(x, y)
  const c = sumMines(n)
  if (c === 0) {
    n.forEach((q) => {
      if (!arena[q[1]][q[0]].open) {
        open(q[0], q[1])
      }
    })
  }
}

function neighbours(x, y) {
  const n = []
  for (let j = y - 1; j <= y + 1; j++) {
    if (j < 0 || j >= arena.length) {
      continue
    }
    const p = arena[j]
    for (let i = x - 1; i <= x + 1; i++) {
      if (i === x && j === y) {
        continue
      }
      if (i < 0 || i >= p.length) {
        continue
      }
      n.push([i, j])
    }
  }
  return n
}

function sumMines(n) {
  let c = 0
  n.forEach((q) => {
    if (arena[q[1]][q[0]].mine) {
      c++
    }
  })
  return c
}

function render() { // updates all the UI world: arena and scores
  for (let j = 0; j < arena.length; j++) {
    const p = arena[j]
    for (let i = 0; i < p.length; i++) {
      const q = p[i]
      const e = q.element
      if (!q.open && !q.flag) {
        e.cont.innerText = ''
        e.div.style.backgroundColor = '#555'
        if (debug()) {
          e.div.style.color = '#333'
          if (q.mine) { e.cont.innerText = 'M' }
        }
        continue
      }
      if (q.flag) {
        e.div.style.backgroundColor = '#555'
        e.div.style.color = '#f00'
        e.cont.innerText = 'F'
      } else {
        if (q.mine) { // oh. hackish
          e.div.style.backgroundColor = '#000'
          e.div.style.color = '#f00'
          e.cont.innerText = 'M'
          e.div.title = 'BOOM!'
          continue
        }
        e.div.style.backgroundColor = '#ccc'
        const nb = sumMines(neighbours(i, j))
        e.div.style.color = ['#ddd', '#009', '#060', '#a00', '#909', '#600', '#099', '#000', '#fff'][nb]
        e.cont.innerText = nb
      }
    }
  }
  const v = document.getElementsByClassName('value')
  v[0].innerText = gameState.lines
  v[1].innerText = gameState.opens
  v[2].innerText = gameState.booms
  persistSave()
}

function delayed(timeout, f) {
  const a = Array.prototype.slice.call(arguments, 2)
  const g = function() {
    f.call(undefined, ...a)
  }
  setTimeout(g, timeout)
}

function initCell(o, table, i, j) { // side effects!
  const cdiv = document.createElement('div')
  cdiv.style.top = `calc(${j} * 4vmin)` // spaces are important
  cdiv.style.left = `calc(${i} * 4vmin)`
  cdiv.addEventListener('click', handler(o, false), false)
  cdiv.addEventListener('contextmenu', handler(o, true), false)
  const cspan = document.createElement('span')
  cdiv.appendChild(cspan)
  table.appendChild(cdiv)
  o.i = i
  o.j = j
  o.element = { div: cdiv, cont: cspan }
}

function newCellObj() {
  return {
    mine: Math.random() < .1,
    open: false,
    flag: false,
  }
}

function initLine() {
  const table = document.getElementById('game')
  const b = []
  const w = arena[0].length
  for (let i = 0; i < w; i++) {
    const o = newCellObj()
    initCell(o, table, i, 0)
    b.push(o)
  }
  return b
}

function init(w, h) {
  const data = persistRestore(w, h)
  if (!data) {
    log('new data')
    arena.length = 0
    for (let j = 0; j < h; j++) {
      const y = []
      for (let i = 0; i < w; i++) {
        y.push(newCellObj())
      }
      arena.push(y)
    }
  } else {
    log('restored data')
    arena.length = 0
    for (let j = 0; j < h; j++) {
      arena.push(data[j])
    }
  }
  const table = document.getElementById('game')
  for (let j = 0; j < arena.length; j++) {
    const a = arena[j]
    for (let i = 0; i < a.length; i++) {
      initCell(a[i], table, i, j)
    }
  }
  render() // TODO remove it?
}

function persistSave() {
  try {
    const x = []
    for (let j = 0; j < arena.length; j++) {
      const y = []
      const b = arena[j]
      for (let i = 0; i < b.length; i++) {
        const c = b[i]
        y.push({
          m: c.mine ? 1 : 0,
          o: c.open ? 1 : 0,
          f: c.flag ? 1 : 0,
        })
      }
      x.push(y)
    }
    localStorage.setItem('x', JSON.stringify({
      arena: x,
      lines: gameState.lines,
      opens: gameState.opens,
      booms: gameState.booms,
    }))
  } catch (e) {
    log(e)
  }
}

function persistRestore(w, h) {
  try {
    const d = JSON.parse(localStorage.getItem('x'))
    const a = []
    const x = d.arena
    if (x.length !== h) {
      throw new Error(`invalid h: ${x.length}`)
    }
    for (let j = 0; j < x.length; j++) {
      const y = x[j]
      if (y.length !== w) {
        return new Error(`invalid w: ${y.length} (j=${j})`)
      }
      const b = []
      for (let i = 0; i < y.length; i++) {
        const q = y[i]
        b.push({
          mine: q.m > 0,
          open: q.o > 0,
          flag: q.f > 0,
        })
      }
      a.push(b)
    }
    gameState.lines = d.lines || 0 // ugly side effect
    gameState.opens = d.opens || 0 // ugly side effect
    gameState.booms = d.booms || 0 // ugly side effect
    return a
  } catch (e) {
    log(e)
  }
  return undefined
}

init(16, 20)
document.getElementById('reset').onclick = () => { // ugly way to restart
  localStorage.clear()
  location.reload() // eslint-disable-line no-undef
}
log('start', Date())
