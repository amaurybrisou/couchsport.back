<template>
  <v-snackbar
    v-model="state"
    :bottom="y === 'bottom'"
    :left="x === 'left'"
    :multi-line="mode === 'multi-line'"
    :right="x === 'right'"
    :top="y === 'top'"
    :vertical="mode === 'vertical'"
    :timeout="3000"
  >
    {{ text }}
    <v-btn color="warning" flat @click="$emit('snack-close')">
      {{ $t('close') }}
    </v-btn>
  </v-snackbar>
</template>

<script>
  export default {
    name: 'AppSnackBar',
    props: {
      state: { type: Boolean, default: false },
      y: { type: String, default: 'top' },
      x: { type: String, default: null },
      mode: { type: String, default: '' },
      text: {
        type: String,
        default: 'Your profile has been successfully saved'
      }
    }
  }
</script>
