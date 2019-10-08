<template>
  <div id="app">
 <GChart
      type="ColumnChart"
      :data="chartData"
      :options="chartOptions"
    />   
  </div>
</template>

<script>
import HelloWorld from './components/HelloWorld.vue'
import { GChart } from 'vue-google-charts';

export default {
  name: 'app',
  components: {
    GChart
  },
  data() {
    return {
      chartData: [
        ['Year', 'Sales', 'Expenses', 'Profit'],
        ['2014', 1000, 400, 200],
        ['2015', 1170, 460, 250],
        ['2016', 660, 1120, 300],
        ['2017', 1030, 540, 350]
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
