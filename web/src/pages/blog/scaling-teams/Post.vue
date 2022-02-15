<template>
  <BlogPost
    title="Scaling teams like parallel computing systems"
    subtitle="Why output is not linear with size"
    date="November 29, 2021"
    :author="author"
    :dive-in-banner="false"
    description="Why team output doesn't scale linearly with team size? Teams are no different from a parallel computing systems. As such, they are subject to Amdahl's law, which governs their scalability."
    reading-time="11 minutes"
    :image="ogImageFull"
  >
    <template #default>
      <p>
        The dynamics of scaling teams and organizations fascinate me. It is well known that the
        productivity of a team does not scale linearly with team size. Going from a team of 5 to a
        team of 10 does not double the output, despite doubling of the team size. The goal of this
        post is to answer <b>why</b> that is, using math.
      </p>
      <p>
        There are of course entire management books dedicated to the subject as well as a lot of
        anecdotal evidence to this claim.
      </p>

      <p>
        A well known rule of thumb is sticking to
        <a
          href="https://docs.aws.amazon.com/whitepapers/latest/introduction-devops-aws/two-pizza-teams.html"
          >Two-Pizza teams</a
        >. In a nutshell, this means a team that can be fed with two pizzas. Personally, I would say
        this <i>really</i> depends on how large the pizzas are but nevertheless, the spirit of the
        message is - there is a productivity sweet spot when it comes to team size. But why is it
        this way?
      </p>

      <p>In this article I make the following conjecture:</p>

      <blockquote>
        A team is no different from a parallel computing system. As such, it is subject to Amdahl's
        law, which governs it's scalability.
      </blockquote>

      <p>
        I am motivated by creating models that approximate different aspects of life. Sometimes such
        models allow us to develop a deeper understanding of the world.
      </p>

      <p>Let's start by dissecting the claim:</p>
      <blockquote>A team is no different from a parallel computing system.</blockquote>

      <h2>Teams <i>are</i> Parallel Computing Systems</h2>

      <h3>Teams are systems</h3>

      <p>
        Per definition, a system is a collection of entities that interact or collaborate in order
        to achieve a common goal. The parts of a system are united in that they have a common
        purpose.
      </p>

      <img
        src="./system.svg"
        class="object-contain w-1/2 sm:w-1/3 mx-auto"
        alt="turing machine illustration"
      />
      <p>
        Similarly, teams are formed in order to achieve objectives that are too big for an
        individual to accomplish - in other words, a shared goal.
      </p>
      <p>
        A common characteristic between systems and teams is the need for interaction and
        collaboration. Without such interaction you have just a collection of workers or entities
        that produce value individually.
      </p>

      <h3>Teams compute</h3>
      <p>Here is an interesting fact:</p>
      <blockquote>The term "Computer" used to refer to a profession!</blockquote>
      <p>
        Fundamentally, computation / calculation simply means transforming one or more inputs into
        one or more outputs. As a team, you outputs may be products or services while the inputs may
        include business requests or your own inspiration.
      </p>

      <p>
        Again, it is key to highlight the need of a common objective in the process of computing -
        in the absence of an aligned common objective, team members are better considered discrete
        systems (teams).
      </p>

      <h3>Teams are (usually) parallel</h3>

      <p>
        Parallel computation is one where multiple calculations are carried out simultaneously. This
        typically applies when problems can be broken down into smaller ones. A classical example of
        a problem well suited for parallelization is
        <a href="https://en.wikipedia.org/wiki/Matrix_multiplication_algorithm"
          >matrix multiplication</a
        >, where Divide-and-Conquer techniques are appropriate.
      </p>
      <p>
        When it comes to teamwork, the system can choose to work on the same task in order to
        improve the quality of the output. For example, this is the case with brainstorming or
        mob/pair programming.
      </p>
      <p>
        More typically however, teams break tasks into smaller ones (Divide-and-Conquer) in order to
        be able to execute on them simultaneously. The sub-tasks are of course related and sometimes
        interdependent because they contribute to the common objective of the system / team.
      </p>

      <p>
        Next, let's discuss parallel systems in general and the Amdahl's law. Later in the article,
        I discuss how we can use this theory in the context of teams.
      </p>

      <h2>Amdahl's scalability law</h2>

      <p>
        In computer systems with more than one processing units, there exists an equation that
        governs the theoretical speedup in execution of a fixed workload as we increase the number
        of processors - known as
        <a href="https://en.wikipedia.org/wiki/Amdahl%27s_law">Amdahl's law</a>.
      </p>
      <p>
        At it's core, this rule stipulates that because workloads always include sections which
        cannot be parallelised (eg. requiring synchronisation between processors), the speedup is
        constrained by the proportion of work that must be performed sequentially.
      </p>
      <p>
        For example if half of the workload cannot be parallelised, then at most we could reduce the
        execution time in half - a speedup of 2x, no matter how much we increase the processor
        count.
      </p>

      <div id="observablehq-workload-b4d16c02" class="text-black"></div>
      <div id="observablehq-viewof-paralleilizable-b4d16c02" class="mt-4"></div>
      <div id="observablehq-best_possible_speedup-20ab6372"></div>

      <h3>Introducing Amdahl's law</h3>

      Amdahl's law is an equation that gives us the theoretical speedup of a workload and is
      formalised as follows:

      <div class="flex">
        <div class="mx-auto text-black text-sm" v-html="amdahl"></div>
      </div>

      <ul>
        where:
        <li>
          <b>p</b> is the proportion of the workload that benefits from parallelization (eg. does
          not require synchronising)
        </li>
        <li><b>n</b> is the parallelism factor (or number of processors)</li>
      </ul>

      <p>
        Let's take for example a workload where 25% of it can not parallelized - this means
        <b>p</b> = 0.75. We can calculate the speedup with parallelism factor (number of processors)
        <b>n</b> = 16. The speedup is just 3.367x, despite 16 processors!
      </p>

      <div class="flex">
        <div class="mx-auto text-black text-sm" v-html="amdahlExample"></div>
      </div>

      <p>
        Even more interesting is plotting the results for <b>p</b> = 0.75. Below, on the x axis we
        have the number of processors and on the y axis we have the resulting speedup. You can
        adjust the <b>p</b> value interactively.
      </p>

      <div id="observablehq-amhdahlPlot-59d50752"></div>
      <div id="observablehq-viewof-p-59d50752"></div>

      <p>
        This visualises the way as parallelism increases, the marginal gain in speedup decreases.
        With a
        <b>p</b> value of 0.75 (meaning 25% of the workload can not be parallelized), as
        <b>n</b> approaches infinity, the speedup approaches 4.
      </p>

      <p>
        More formally, we can say that as <b>n</b> grows towards infinity, the Speedup tends to
        1/(1-p).
      </p>

      <p>
        A key observation to be made here is that the slope of the curve platoes. In other words,
        the marginal effect of increasing n diminishes. Let's look further into that.
      </p>

      <h3>Rate of speedup change</h3>

      <p>
        The first order derivative (rate of change) is the curve slope of our function. This means
        it gives us an indication of the sensitivity with respect to the parallelism factor
        <b>n</b>. A value of 1 would indicate linear relationship, and we can see that as n
        increases, the rate of change tends towards 0.
      </p>

      <div class="flex">
        <div class="mx-auto text-black text-sm" v-html="rate"></div>
      </div>

      <p>
        This is the formal definition of the first order derivative of the Amdahl's law. Let's plot
        it!
      </p>

      <div id="observablehq-rateOfChange-af23ee38"></div>
      <div id="observablehq-viewof-p2-af23ee38"></div>

      <p>
        It's clear that even with high values of <b>p</b>, there is a quick dropoff in sensitivity
        towards <b>n</b>. This insight is useful because n is a finite resource - processor units.
        In this context <b>f'(n)</b> represents a measure of efficiency. A greater rate of change
        for increasing <b>n</b> means higher efficiency.
      </p>

      <p>
        Returning to the example where three quarters of the workload are parallelizable (p = 0.75),
        an argument could be made that it is wasteful to allocate more than 16 processor units (n =
        16) and depending on exactly how scarce the resource is, perhaps n = 10 is a more
        appropriate allocation.
      </p>

      <p>
        However, if the workload could be modified to have a greater parallelizable proportion
        <b>p</b>, parallelization could be efficient at higher <b>n</b> values!
      </p>

      <p>This brings up an interesting question:</p>
      <blockquote>
        If we can establish a threshold of a minimum acceptable efficiency f'(n), how much of the
        workload needs to be parallelizable to efficiently utilize a given number of processors?
      </blockquote>
      <p>We can answer this by rearranging the equation and solve for <b>p</b>.</p>

      <h3>Solving for p</h3>

      <p>
        So far we had the parallel proportion p set and we plotted the Speedup and the rate of
        change with respect of the parallelism factor n.
      </p>

      <div class="flex">
        <div class="mx-auto text-black text-sm" v-html="forp"></div>
      </div>

      <p>
        Since we use the rate of change as our efficiency measure, let's solve for p instead. Let's
        refer to the rate of change <b>f'(n)</b> as <b>d</b> for brevity. The equation is starting
        to look bulky but we can plot it and reason around it visually.
      </p>

      <div id="observablehq-forp-54b848a7"></div>
      <div id="observablehq-viewof-d-54b848a7"></div>

      <p>
        If we can establish a reasonable minimum efficiency <b>d</b> for a given workload, we can
        see what proportion <b>p</b> of that workload must be parallelizable for different
        <b>n</b> values.
      </p>

      <p>
        The beautiful part is that this equation gives us an indication as to how the workload must
        be adjusted in order to be able to continue scaling! In other words:
      </p>
      <blockquote>
        This allows us to reason around the inherent bottleneck caused by the need of task
        synchronization.
      </blockquote>

      <p>
        If we choose a minimum efficiency <b>d</b> (or worst rate of improvement per processor
        added) of 0.2, and if we want to utilize 5 processors, then 76% of the workload must be
        parallelizable. On the other hand, if we want to utilize 20 processors and maintain the same
        level of efficiency, then 94% of the workload needs to be parallelizable.
      </p>

      <p>Let's get back to discussing teams.</p>

      <h2>Scaling teams</h2>

      <p>
        So far I have made the case that teams are parallel computing systems. We also discussed the
        Amdahl's law, which governs the speedup in completing a workload (objective) gained from
        increasing the number of processing units.
      </p>

      <p>
        When we made the case that teams are parallel computing systems, we put an emphasis of a
        common goal / objective of the team / system. There is an inherent need for interaction and
        collaboration between the parts of the system. Without it you have just a collection of
        workers or entities that produce value individually. Building on this:
      </p>

      <blockquote>
        The common team goal or objective is its workload. Parts of the workload (subtasks) may be
        dependent on each other, and some subtasks may require synchronization between team members.
      </blockquote>

      <p>
        Let's consider the points of collaboration within a team as serial (non-parallelizable)
        portions of their workload. For example it is necessary for a healthy team to:
      </p>
      <ul>
        <li>coordinate work (planning)</li>
        <li>share updates (standups)</li>
        <li>re-evaluate & optimize their process (retrospectives)</li>
      </ul>

      <p>
        And now the <b>really cool part</b> - If we take a look at how a team spends their week, we
        can approximate the proportion of the workload that is non-parallelizable. Why is this
        useful? We can use the Amdahl's law and it's first order derivative to reason around team
        dynamics. For instance:
      </p>
      <blockquote>
        We can answer the question "What would be the maximum team size where team members can feel
        reasonably productive?" from first principles.
      </blockquote>

      <p>
        Let's do it! First, we assume a 40 hour working week. How many of those hours are spent for
        the healthy operation of a team collaborating towards a common objective? This is a portion
        of the work that can not be parallelized and for a good reason. This number can vary, but 5
        hours would be reasonable, which represents 12.5% of the total time. In other words, the
        parallelizable proportion <b>p</b> is 0.875.
      </p>

      <p>
        Our second question is - what is a reasonable minimum efficiency in a team setting? Earlier,
        when discussing the rate of change with respect to the number of processors n, a value d of
        0.1 may have been reasonable - CPUs cores are cheap. For a team a much higher value is
        desirable, I would say 0.3 or higher.
      </p>

      <p>A <b>p</b> = 0.875 and <b>d</b> = 0.3 gives us a team size <b>n</b> of 6.68.</p>

      <h3>Solving for p</h3>

      <p>
        As earlier, let's plot this, solving for <b>p</b>. We can pick a minimum efficiency value
        <b>d</b> and solve for what ratio of the workload must be parallelizable. For convenience,
        below I am also converting <b>p</b> to the number of hours per week that can be spent in
        synchronization tasks.
      </p>

      <div id="observablehq-serilHours-af307419"></div>
      <div id="observablehq-viewof-d2-af307419"></div>

      <p>
        If you have in mind the number of hours your team needs to spend synchronizing every week,
        this plot visualizes how large the given team can be!
      </p>
      <p>
        Keep in mind the meaning of the parameter <b>d</b> here - it is a factor that determines the
        marginal increase in speedup/output when adding an additional processing unit <b>n</b>.
      </p>

      <h2>Conclusions</h2>

      <blockquote>Teams are parallel computer systems who eat pizzas.</blockquote>
      <p>
        Because by definition a team works towards a common goal, there is an inherent need for
        synchronization between the members of the team. 12% is a reasonable estimate of the
        proportion of time spend synchronizing.
      </p>

      <p>
        The proportion of time that cannot be parallelized places an upper bound of how much of the
        workload can be sped up. Moreover, the rate of change in output decreases as the number of
        team member increases. In other words:
      </p>

      <blockquote>
        After a certain threshold, increasing a team size is a very inefficient way of reducing the
        time to achieve the team's goal (boosting output).
      </blockquote>

      <p>
        So, what is the right team size? Assuming 12% of the workload needs to be synchronized to
        facilitate healthy collaboration within the team, and if we want a reasonable minimum
        efficiency of 0.3, then the answer is a size of 6-7. This happens to be align really well
        with the "Two-Pizzas Team" size!
      </p>

      <p>In my next post, I will control for different pizza sizes!</p>

      <p>
        Thanks for reading!
        <br />
        - <a href="https://twitter.com/krlvi">Kiril</a>
      </p>
      <hr />
    </template>
  </BlogPost>
</template>

<script lang="ts" setup>
import BlogPost from '../BlogPost.vue'
import avatar from '../kiril.jpeg'
import { Inspector, Runtime } from '@observablehq/runtime'
import katex from 'katex'
import 'katex/dist/katex.min.css'
import ogImage from './system.png'
import { computed, onMounted } from 'vue'

const author = {
  name: 'Kiril Videlov',
  avatar: avatar,
  link: 'https://twitter.com/krlvi',
}

let ogImageFull = computed(() => new URL(ogImage, location.origin).href)

let amdahl = computed(() =>
  katex.renderToString(String.raw`Speedup(n) = \frac{1}{(1-p) + \frac{p}{n}}`, {
    displayMode: true,
    throwOnError: false,
  })
)

let amdahlExample = computed(() =>
  katex.renderToString(String.raw`\frac{1}{(1-0.75) + \frac{0.75}{16}}  = 3.367`, {
    displayMode: true,
    throwOnError: false,
  })
)

let rate = computed(() =>
  katex.renderToString(String.raw`f'(n) = \frac{p}{((p-1)n-p)^2}`, {
    displayMode: true,
    throwOnError: false,
  })
)

let forp = computed(() =>
  katex.renderToString(
    String.raw`p = \frac{-(\sqrt{4dn^2 - 4dn +1} -2dn^2 + 2dn -1)}{2dn^2 - 4dn +2d}`,
    {
      displayMode: true,
      throwOnError: false,
    }
  )
)

onMounted(async () => {
  let module = await import('https://api.observablehq.com/d/1c779005074f8668.js?v=3')
  new Runtime().module(module.default, (name) => {
    if (name === 'workload')
      return new Inspector(document.querySelector('#observablehq-workload-b4d16c02'))
    if (name === 'viewof paralleilizable')
      return new Inspector(document.querySelector('#observablehq-viewof-paralleilizable-b4d16c02'))
    if (name === 'best_possible_speedup')
      return new Inspector(document.querySelector('#observablehq-best_possible_speedup-20ab6372'))
    if (name === 'amhdahlPlot')
      return new Inspector(document.querySelector('#observablehq-amhdahlPlot-59d50752'))
    if (name === 'viewof p')
      return new Inspector(document.querySelector('#observablehq-viewof-p-59d50752'))
    if (name === 'rateOfChange')
      return new Inspector(document.querySelector('#observablehq-rateOfChange-af23ee38'))
    if (name === 'viewof p2')
      return new Inspector(document.querySelector('#observablehq-viewof-p2-af23ee38'))
    if (name === 'forp') return new Inspector(document.querySelector('#observablehq-forp-54b848a7'))
    if (name === 'viewof d')
      return new Inspector(document.querySelector('#observablehq-viewof-d-54b848a7'))
    if (name === 'serilHours')
      return new Inspector(document.querySelector('#observablehq-serilHours-af307419'))
    if (name === 'viewof d2')
      return new Inspector(document.querySelector('#observablehq-viewof-d2-af307419'))
  })
})
</script>
