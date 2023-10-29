/*global document:readable, setTimeout:readable*/
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
  let t = ''
  for (let i = 0; i < arguments.length; i++) {
    if (t !== '') {
      t += ' '
    }
    t += arguments[i]
  }
  const h = document.getElementById('debug')
  const e = document.createElement('div')
  e.innerText = t
  h.prepend(e)
}

const arena = [] // it has to be split: data (mutable), pointers to DOM (immutable)
let locked = false

function handler(obj, rightClick) {
  return (ev) => {
    ev.stopPropagation()
    ev.preventDefault()
    if (locked) {
      log('LOCKED')
      return
    }
    const x = obj.i
    const y = obj.j
    log('click', x, y, rightClick)
    if (rightClick) {
      log(arena[y][x].flag, !arena[y][x].flag)
      arena[y][x].flag = !arena[y][x].flag
    } else {
      open(x, y)
      collapse() // ugly: collapse only if opened; return if no more collapsing
    }
    render() // ugly: update entire arena
    return false
  }
}

function collapse() {
  locked = true
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
  locked = false // unlock if nothing more to do
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
      arena[j][i].element.div.style.top = `calc(${v} * 4vw)` // spaces are important
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
  arena.unshift(initLine(document.getElementById('game'), 0, arena[0].length)) // getElementById — ugly; it has to be factory
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
      if (q.mine) {
        e.cont.innerText = 'M'
      } else {
        e.cont.innerText = sumMines(neighbours(i, j))
      }
      e.div.style.backgroundColor = q.open ? '#fff' : '#ccc' // TODO debug
      e.cont.style.color = (q.mine && !q.flag) ? '#fff' : (q.flag ? '#c00' : '#000') // TODO debug
    }
  }
}

function delayed(timeout, f) {
  const a = Array.prototype.slice.call(arguments, 2)
  const g = function() {
    f.call(undefined, ...a)
  }
  setTimeout(g, timeout)
}

function initLine(table, j, w) {
  const b = []
  for (let i = 0; i < w; i++) {
    const o = {}
    const cdiv = document.createElement('div')
    cdiv.style.top = `calc(${j} * 4vw)` // spaces are important
    cdiv.style.left = `calc(${i} * 4vw)`
    cdiv.addEventListener('click', handler(o, false), false)
    cdiv.addEventListener('contextmenu', handler(o, true), false)
    const cspan = document.createElement('span')
    cdiv.appendChild(cspan)
    table.appendChild(cdiv)
    o.i = i // ugly: it has to be factory of objects with convenient interface or functions to manipulate with
    o.j = j
    o.mine = Math.random() < .1
    o.element = { div: cdiv, cont: cspan }
    o.open = false
    o.flag = false
    b.push(o)
  }
  return b
}

function init(w, h) {
  arena.length = 0
  const table = document.getElementById('game')
  for (let j = 0; j < h; j++) {
    arena.push(initLine(table, j, w))
  }
  render() // TODO remove it
}

init(16, 20)
log('start', Date())
