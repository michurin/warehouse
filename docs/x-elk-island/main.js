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

function log() {
  const t = []
  const s = []
  for (let i = 0; i < arguments.length; i++) {
    t.push(arguments[i])
    if (arguments[i].stack) {
      s.push(arguments[i].stack)
    }
  }
  const h = document.getElementById('debug')
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
      evaporate(j)
      return // without unlocking
    }
  }
  unlock() // unlock if nothing more to do
}

function evaporate(k) {
  let tmp = k // TODO REMOVE it
  for (let i = 0; i < arena[k].length; i++) {
    arena[k][i].element.div.style.borderColor = '#0f0'
    tmp = `${tmp} /${arena[k][i].open ? 'O' : ''}${arena[k][i].mine ? 'M' : ''}`
  }
  log('evaporate', tmp)
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
  if (step < 10) {
    delayed(20, evaporate_remove, k, step + 1)
    return
  }
  for (let j = 0; j < k; j++) {
    for (let i = 0; i < arena[j].length; i++) {
      arena[j][i].j++
    }
  }
  arena.splice(k, 1)
  arena.unshift(initLine())
  render()
  delayed(500, evaporate_open)
}

function evaporate_open() {
  for (let j = 0; j < arena.length; j++) {
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

function render() {
  for (let j = 0; j < arena.length; j++) {
    const p = arena[j]
    for (let i = 0; i < p.length; i++) {
      const q = p[i]
      const e = q.element
      if (!q.open && !q.flag) {
        e.cont.innerText = ''
        e.div.style.backgroundColor = '#555'
        // hack
        e.div.style.color = '#333'
        if (q.mine) { e.cont.innerText = 'M' }
        // /hack
        continue
      }
      if (q.flag) {
        e.div.style.backgroundColor = '#555'
        e.div.style.color = '#f00'
        e.cont.innerText = 'F'
      } else {
        e.div.style.backgroundColor = '#ccc'
        const nb = sumMines(neighbours(i, j))
        e.div.style.color = ['#ddd', '#009', '#060', '#a00', '#909', '#600', '#099', '#000', '#fff'][nb]
        e.cont.innerText = nb
      }
    }
  }
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
    localStorage.setItem('x', JSON.stringify({ x: x }))
  } catch (e) {
    log(e)
  }
}

function persistRestore(w, h) {
  try {
    const d = JSON.parse(localStorage.getItem('x'))
    const a = []
    const x = d.x
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
    return a
  } catch (e) {
    log(e)
  }
  return undefined
}

init(16, 20)
log('start', Date())
