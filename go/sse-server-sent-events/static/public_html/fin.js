'use strict'

const params = new URLSearchParams(window.location.search)
const room = params.get('back')

const r = room.replaceAll(/[^0-9a-zA-Z_-]+/g, '')
document.getElementsByTagName('b')[0].innerText = r
document.getElementsByTagName('a')[0].href = '/' + r

const reason = params.get('reason')
const k = reason.replaceAll(/[^a-z]+/g, '')
document.getElementsByTagName('h2')[0].innerText = k // TODO key for translation

