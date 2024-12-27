<script lang="ts">
	import Panel from "$lib/panel.svelte";
	import { currentElection, elections } from "../store";
	import { goto } from "$app/navigation";

	$: upcomingElections = $elections?.filter((e) => !e.isActive && !e.isConcluded) ?? [];

	$: if ($currentElection && !$currentElection.hasVoted) {
		goto(`/vote`);
	} else if (upcomingElections.length > 0) {
		goto(`/election/${upcomingElections[0].id}`);
	}
</script>

<svelte:head>
	<title>CSS Elects</title>
</svelte:head>

<Panel title="There are no upcoming elections">
	<p>Check this space later for updates.</p>
	<img
		src={`https://www.guildofstudents.com/asset/Organisation/6419/cathsoc%20logo.png`}
		height="100px"
	/>
</Panel>

<style>
	img {
		margin-top: 16px;
	}
</style>
