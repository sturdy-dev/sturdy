<template>
  <Documentation>
    <JobDescription title="Senior Backend Software Engineer">
      <p>Hey! We’re <strong>Sturdy – simple & efficient code collaboration.</strong></p>

      <p>
        We where part of Y Combinator in 2021, and have raised a seed round from several tech funds,
        and a team of amazing angel investors.
      </p>

      <p>We’re now looking to hire more engineers to come and join our team!</p>

      <h2>You’ll work in our team and...</h2>

      <ul>
        <li>... build our open-source backend and CLI (in Go!)</li>
        <li>... enjoy what you’re doing (which is programming!)</li>
      </ul>

      <p>
        Because we’re an early-stage company, we’re looking for a motivated individual who can help
        making Sturdy a better product.
      </p>

      <h2>How we build our backend</h2>

      <p>
        The backend at Sturdy powers the heavy lifting of everything related to files (and we have a
        lot of them!). It's built with Go, PostgreSQL, and GraphQL. Our backend is a single
        monolith, that we're deploying to Sturdy Cloud on AWS, and to self-hostable installations at
        our customers (in Docker). The backend is
        <a href="https://github.com/sturdy-dev/sturdy/tree/main/api">open-source</a> and all of the
        work happens in the open.
      </p>

      <h2>Skills of interest</h2>

      <ul>
        <li>Solid experience and interest to work with Go, GraphQL, and PostgreSQL</li>
        <li>Have a solid understanding of distributed systems</li>
        <li>Flexibility to adjust to the dynamic nature of a startup</li>
      </ul>

      <h2>Together with the team, you’ll succeed with these goals:</h2>

      <ul>
        <li>Work as a part of a team of engineers to add features to the Sturdy service</li>
        <li>Have a good overview of the product, from customer needs, to product features</li>
      </ul>
    </JobDescription>
  </Documentation>
</template>

<script lang="ts" setup>
import JobDescription from '../../organisms/careers/JobDescription.vue'
import Documentation from '../../layouts/Documentation.vue'
</script>
