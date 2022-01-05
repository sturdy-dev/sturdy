<template>
  <StaticPage
    title="Quickstart"
    subtitle="Get up and running with Sturdy and start collaborating on code without Git."
    category="documentation"
    metadescription="Getting started with Sturdy"
  >
    <div class="mt-5 text-gray-500 mx-auto lg:max-w-none lg:row-start-1 lg:col-start-1">
      <div class="prose prose-yellow">
        <p>
          Coding on Sturdy can begin as soon as you create a Sturdy account, and requires three
          steps:
        </p>

        <ol>
          <li>
            Install the client app so that your project files can be made available in a folder on
            your computer
          </li>
          <li>Create a blank codebase or import a Git repository</li>
          <li>Connect the client app to your newly created codebase</li>
        </ol>

        <hr />

        <h2 class="flex">1. Install the Sturdy client</h2>
        <p>
          Sturdy allows real-time code collaboration and feedback for you and your team. This does
          not mean you need to code in the browser! The Sturdy sets up a directory on your computer
          that will be in sync with the cloud.
        </p>
      </div>

      <InstallationInstructions class="border-l-2 border-black pl-4 py-0 my-4" />

      <div class="prose prose-yellow">
        <h2 class="flex">2. Create a blank codebase or import an existing one</h2>

        <p>
          A codebase in Sturdy is just a collection of files that belong together, similar to a Git
          repository. In fact, Sturdy allows you to import an existing Git repository as well.
        </p>

        <div class="relative text-base mx-auto max-w-prose lg:max-w-none">
          <figure>
            <div class="lg:aspect-w-12 aspect-h-7 lg:aspect-none">
              <img
                class="rounded-lg object-cover object-center"
                alt="Creating a Sturdy codebase"
                src="../../assets/Web/CreateCodebase.png"
                height="658"
                width="801"
              />
            </div>
          </figure>
        </div>

        <h2 class="flex">3. Connect the client app</h2>
        <p>
          After creating or importing your code into Sturdy you will see a prompt to connect the
          codebase to your computer.
        </p>
        <div class="relative text-base mx-auto max-w-prose lg:max-w-none">
          <figure>
            <div class="lg:aspect-w-12 aspect-h-7 lg:aspect-none">
              <img
                class="rounded-lg object-cover object-center"
                alt="Connecting a Sturdy codebase to a computer"
                src="../../assets/Web/ConnectCodebaseToComputer.png"
                height="214"
                width="782"
              />
            </div>
            <figcaption class="mt-3 flex text-sm text-gray-500">
              <span class="ml-2"
                >Connecting a Sturdy codebase to a computer. Note that the codebase name and ID will
                be different for you.</span
              >
            </figcaption>
          </figure>
        </div>

        <p>
          Copy and run the provided command in your terminal &mdash; this will initialize the
          codebase in a folder on your computer.
        </p>
        <p>
          The last argument (in this example: <code>my-first-codebase</code>) is the path to the
          directory that Sturdy will create and setup your project in. If you want to setup Sturdy
          in another directory, change the last argument with the new path.
        </p>
        <p>
          If this is the first codebase that you run <code>sturdy init</code>, you will be asked to
          authenticate yourself. Follow the on screen instructions.
        </p>
        <div class="relative text-base mx-auto max-w-prose lg:max-w-none">
          <figure>
            <div class="lg:aspect-w-12 aspect-h-7 lg:aspect-none">
              <img
                class="rounded-lg object-cover object-center"
                alt="Authenticating the Sturdy client"
                src="../../assets/Web/AuthenticateClient.png"
                height="98"
                width="782"
              />
            </div>
            <figcaption class="mt-3 flex text-sm text-gray-500">
              <span class="ml-2">Authenticating the Sturdy client on first use.</span>
            </figcaption>
          </figure>
        </div>

        <p>After authenticating the client sets up the codebase in a folder on your computer.</p>

        <div class="relative text-base mx-auto max-w-prose lg:max-w-none">
          <figure>
            <div class="lg:aspect-w-12 aspect-h-7 lg:aspect-none">
              <img
                class="rounded-lg object-cover object-center"
                alt="Successful Sturdy client setup"
                src="../../assets/Web/AuthenticateClientSuccess.png"
              />
            </div>
            <figcaption class="mt-3 flex text-sm text-gray-500">
              <span class="ml-2"
                >Setup is successful. Note that the names and IDs will be different for you.</span
              >
            </figcaption>
          </figure>
        </div>

        <h2>Next steps</h2>
        <p>
          After creating a codebase and connecting it to your computer you are ready to start
          coding. If you want to read more before getting started, here is
          <router-link :to="{ name: 'resourcesDocs' }">more of our documentation</router-link>.
        </p>
      </div>
    </div>
  </StaticPage>
</template>

<script>
import StaticPage from './StaticPage.vue'
import InstallationInstructions from '../../components/install/InstallationInstructions.vue'
export default {
  components: { InstallationInstructions, StaticPage },
}
</script>
