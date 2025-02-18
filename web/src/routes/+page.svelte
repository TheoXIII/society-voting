<script lang="ts">
	import Panel from "$lib/panel.svelte";
	import { currentPoll, polls } from "../store";
	import { goto } from "$app/navigation";

	$: upcomingPolls = $polls?.filter((e) => !e.isActive && !e.isConcluded) ?? [];

	$: if ($currentPoll && !$currentPoll.hasVoted) {
		goto(`/vote`);
	} else if (upcomingPolls.length > 0) {
		let poll = upcomingPolls[0];
		goto(`/${poll.pollType.name.toLowerCase()}/${poll.id}`);
	}
</script>

<svelte:head>
	<title>CathSoc Elects</title>
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
