import Vue from 'vue'
import Router from 'vue-router'
import Home from '@/views/Home'
import Monitor from '@/views/Monitor'

Vue.use(Router)

const NotFoundView = Vue.component('NotFoundView', {
  template: '<h1>...Ops, 404</h1>'
})

export default new Router({
  // mode: 'history',
  routes: [
      {
      path: '/',
      redirect: '/home'
      },
      {
          path: '/home',
          name: 'Home',
          component: Home
      },
      {
        path: '/monitor',
        name: 'monitor',
        component: Monitor
    }
  ]
})
