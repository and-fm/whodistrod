<script lang="ts">
  import { AppConfig } from "$lib/config";

  let trackUrl = "";
  let distributor: string | null = null;
  let loading = false;
  let error: string | null = null;

  function parseTrackId(url: string): string | null {
    try {
      const parsedUrl = new URL(url);

      // Handle /browse/track/<id> format
      if (parsedUrl.pathname.includes("/browse/track/")) {
        const match = parsedUrl.pathname.match(/\/browse\/track\/(\d+)/);
        return match?.[1] || null;
      }

      // Handle /album/<albumId>/track/<trackId> format
      if (parsedUrl.pathname.includes("/album/")) {
        const match = parsedUrl.pathname.match(/\/track\/(\d+)/);
        return match?.[1] || null;
      }

      // Handle direct /track/<id> format (if it exists)
      if (parsedUrl.pathname.includes("/track/")) {
        const match = parsedUrl.pathname.match(/\/track\/(\d+)/);
        return match?.[1] || null;
      }

      return null;
    } catch (e) {
      return null;
    }
  }

  async function handleSubmit() {
    loading = true;
    error = null;
    try {
      const trackId = parseTrackId(trackUrl);
      if (!trackId) {
        throw new Error(
          "Invalid Tidal Track URL. Please check the URL and try again."
        );
      }

      const response = await fetch(
        `${AppConfig.backendUrl}/v1/tidal/track/${trackId}/providers`,
        {
          method: "GET",
          headers: {
            "Content-Type": "application/json",
          },
        }
      );

      if (!response.ok) {
        throw new Error("Failed to fetch distributor information");
      }

      const data = await response.json();
      distributor = data.provider;
    } catch (e) {
      error = e.message;
      resetState();
    } finally {
      loading = false;
    }
  }

  function resetState() {
    distributor = null;
  }
</script>

<div
  class="min-h-screen bg-white text-black flex flex-col items-center px-4 py-12 md:py-24"
>
  <div class="w-full max-w-2xl mx-auto">
    <div class="space-y-8 mb-16">
      <h1 class="text-5xl md:text-6xl font-normal tracking-tight">
        who distro'd?
      </h1>
      <p class="text-lg md:text-xl text-gray-600 max-w-xl">
        instantly discover which distributor delivered a given Tidal track
      </p>
    </div>

    <form on:submit|preventDefault={handleSubmit} class="space-y-6 mb-12">
      <div class="relative">
        <input
          type="url"
          bind:value={trackUrl}
          placeholder="paste your tidal track url here"
          class="w-full px-6 py-4 bg-transparent border-2 border-black rounded-none focus:outline-none focus:border-blue-600 text-lg placeholder:text-gray-400"
          required
        />
      </div>

      <button
        type="submit"
        class="w-full md:w-auto px-8 py-4 bg-black text-white hover:bg-gray-800 transition-colors duration-200 text-lg cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
        disabled={loading}
      >
        {#if loading}
          searching...
        {:else}
          find distributor
        {/if}
      </button>
    </form>

    {#if error}
      <div class="mb-8 py-4 text-red-600 text-lg">
        {error}
      </div>
    {/if}

    {#if distributor && !error}
      <div class="bg-emerald-400 border-2 border-black p-8">
        <div class="text-lg text-black/70 uppercase tracking-wider mb-2">
          track distributor
        </div>
        <h2 class="text-4xl md:text-5xl font-medium text-black mb-4">
          {distributor}
        </h2>
        <div class="w-12 h-0.5 bg-black/40"></div>
      </div>
    {/if}
  </div>

  <footer
    class="fixed bottom-0 left-0 w-full py-6 px-4 bg-white border-t border-gray-200"
  >
    <div
      class="max-w-2xl mx-auto flex justify-between items-center text-sm text-gray-600"
    >
      <div>© 2025 and.fm</div>
      <div class="flex space-x-6">
        <a href="https://and.fm" target="_blank" class="hover:text-black"
          >and.fm</a
        >
        <a
          href="https://github.com/and-fm/whodistrod"
          target="_blank"
          class="hover:text-black">github</a
        >
      </div>
    </div>
  </footer>
</div>
