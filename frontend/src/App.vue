<template>
  <div id="app">
 <GChart
      type="CandlestickChart"
      :data="chartData"
      :options="chartOptions"
    />   
  </div>
</template>

<script>
import { GChart } from 'vue-google-charts';

export default {
  name: 'app',
  components: {
    GChart
  },
  data() {
    return {
      chartData: [
              ['Mon', 20, 28, 38, 45],
              ['Tue', 31, 38, 55, 66],
              ['Wed', 50, 55, 77, 80],
              ['Thu', 77, 77, 66, 50],
              ['Fri', 68, 66, 22, 15]
      ],
      chartOptions: {
        title: 'Company Performance',
        subtitle: 'Sales'
      }
    };
  }
}
</script>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
