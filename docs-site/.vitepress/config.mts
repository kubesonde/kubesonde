import { defineConfig } from 'vitepress'

// Kubesonde documentation site.
// Content is organised with the Diátaxis framework: Tutorials (learning),
// How-to guides (goals), Concepts (understanding) and Reference (information).
export default defineConfig({
  title: 'Kubesonde',
  description:
    'Probe and visualize the actual network connectivity of applications in a Kubernetes cluster, and compare it against the policies you meant to enforce.',
  lang: 'en-US',
  base: '/kubesonde/',
  cleanUrls: true,
  lastUpdated: true,
  // localhost URLs in port-forward examples are intentional, not dead links.
  ignoreDeadLinks: [/^https?:\/\/localhost/],
  head: [
    ['meta', { name: 'theme-color', content: '#0b6b5f' }],
    ['meta', { property: 'og:type', content: 'website' }],
    ['meta', { property: 'og:title', content: 'Kubesonde Documentation' }],
    [
      'meta',
      {
        property: 'og:description',
        content:
          'Diagnose the gap between intended and actual Kubernetes network connectivity.',
      },
    ],
  ],

  themeConfig: {
    logo: '/logo.png',
    siteTitle: 'Kubesonde',

    nav: [
      { text: 'Guide', link: '/guide/introduction', activeMatch: '/guide/' },
      { text: 'Tutorials', link: '/tutorials/first-scan', activeMatch: '/tutorials/' },
      { text: 'How-to', link: '/how-to/install', activeMatch: '/how-to/' },
      { text: 'Concepts', link: '/concepts/how-it-works', activeMatch: '/concepts/' },
      { text: 'Reference', link: '/reference/kubesonde-resource', activeMatch: '/reference/' },
      {
        text: 'v1',
        items: [
          { text: 'Releases', link: 'https://github.com/kubesonde/kubesonde/releases' },
          { text: 'Hosted UI', link: 'https://kubesonde.jackops.dev' },
        ],
      },
    ],

    sidebar: {
      '/guide/': [
        {
          text: 'Getting started',
          items: [
            { text: 'Introduction', link: '/guide/introduction' },
            { text: 'Quickstart', link: '/guide/quickstart' },
            { text: 'Installation', link: '/guide/installation' },
          ],
        },
        {
          text: 'Going further',
          items: [
            { text: 'Tutorials', link: '/tutorials/first-scan' },
            { text: 'How-to guides', link: '/how-to/install' },
            { text: 'Concepts', link: '/concepts/how-it-works' },
            { text: 'Reference', link: '/reference/kubesonde-resource' },
          ],
        },
      ],
      '/tutorials/': [
        {
          text: 'Tutorials',
          items: [
            { text: 'Your first connectivity scan', link: '/tutorials/first-scan' },
            { text: 'Auditing a NetworkPolicy', link: '/tutorials/audit-networkpolicy' },
          ],
        },
      ],
      '/how-to/': [
        {
          text: 'How-to guides',
          items: [
            { text: 'Install & uninstall', link: '/how-to/install' },
            { text: 'Target specific pods', link: '/how-to/target-pods' },
            { text: 'Assert expected outcomes', link: '/how-to/assert-outcomes' },
            { text: 'Fetch & export results', link: '/how-to/fetch-results' },
            { text: 'Visualize results', link: '/how-to/visualize' },
            { text: 'Tune probe concurrency', link: '/how-to/tune-concurrency' },
          ],
        },
      ],
      '/concepts/': [
        {
          text: 'Concepts',
          items: [
            { text: 'How Kubesonde works', link: '/concepts/how-it-works' },
            { text: 'Intended vs. actual connectivity', link: '/concepts/intended-vs-actual' },
            { text: 'Architecture', link: '/concepts/architecture' },
            { text: 'Probe completeness (quiescence)', link: '/concepts/completeness' },
            { text: 'What Kubesonde is not', link: '/concepts/what-it-is-not' },
          ],
        },
      ],
      '/reference/': [
        {
          text: 'Reference',
          items: [
            { text: 'Kubesonde resource', link: '/reference/kubesonde-resource' },
            { text: 'Controller HTTP API', link: '/reference/http-api' },
            { text: 'Environment variables', link: '/reference/environment' },
            { text: 'kubectl cheat sheet', link: '/reference/kubectl-cheatsheet' },
          ],
        },
      ],
    },

    outline: { level: [2, 3], label: 'On this page' },

    search: { provider: 'local' },

    socialLinks: [
      { icon: 'github', link: 'https://github.com/kubesonde/kubesonde' },
    ],

    editLink: {
      pattern: 'https://github.com/kubesonde/kubesonde/edit/main/docs-site/:path',
      text: 'Edit this page on GitHub',
    },

    footer: {
      message: 'Released under the Apache License 2.0.',
      copyright: 'Copyright © 2023–present Kubesonde contributors',
    },

    docFooter: { prev: 'Previous', next: 'Next' },
  },
})
