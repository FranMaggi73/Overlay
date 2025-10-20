<script>
	import { onMount } from 'svelte';
	
	let gifUrl = '/gifs/pet.gif'; // GIF por defecto desde carpeta local
	
	// Detectar si estamos en Electron
	onMount(() => {
		isElectron = navigator.userAgent.toLowerCase().indexOf('electron') > -1;
	});
</script>

<main class="overlay-container">
	<!-- GIF animado en la esquina -->
	<div class="gif-container">
		<img src={gifUrl} alt="Animated GIF" class="overlay-gif" />
	</div>
</main>

<style>
	:global(body) {
		margin: 0;
		padding: 0;
		overflow: hidden;
		font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
	}
	
	.overlay-container {
		width: 100vw;
		height: 100vh;
		background: transparent;
		position: relative;
		pointer-events: none; /* IMPORTANTE: hace que el fondo no bloquee clicks */
	}
	
	
	.gif-container {
		position: absolute;
		bottom: 40px;
		right: 40px;
		pointer-events: none;
	}
	
	.overlay-gif {
		width: 80px;
		height: 80px;
		object-fit: contain;
		border-radius: 12px;
		box-shadow: 0 4px 20px rgba(0, 0, 0, 0.5);
		animation: float 3s ease-in-out infinite;
	}
	
	/* Animaciones */
	@keyframes pulse {
		0%, 100% {
			opacity: 1;
			transform: scale(1);
		}
		50% {
			opacity: 0.5;
			transform: scale(1.3);
		}
	}
	
	@keyframes bounce {
		0%, 100% {
			transform: translateY(0);
		}
		50% {
			transform: translateY(-10px);
		}
	}
	
	@keyframes float {
		0%, 100% {
			transform: translateY(0px);
		}
		50% {
			transform: translateY(-20px);
		}
	}
</style>