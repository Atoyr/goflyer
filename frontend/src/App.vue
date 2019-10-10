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
              ['time', 'open','high','low','close']
      ],
      chartOptions: {
        title: 'Company Performance',
        subtitle: 'Sales',
        bar: { groupWidth: '50%' }, // Remove space between bars.
        candlestick: {
          fallingColor: { strokeWidth: 0, fill: '#a52714' }, // red
          risingColor: { strokeWidth: 0, fill: '#0f9d58' }   // green
        }
      }
    }
  },
  created() {
    fetch('http://localhost:8080/v1/Candlestick/BTC_JPY/3m')
    .then(response => {
      return response.json()
    })
    .then(json => {
      console.log(json)
      for (var i = 0; i < json.length;i++ ){
          console.log(json[i].time)
            console.log(json[i].open)
            console.log(json[i].high)
            console.log(json[i].low)
            console.log(json[i].close )
          var a = new Array(json[i].time,  json[i].high,json[i].open,json[i].close, json[i].low  )
        this.chartData.push(a)
      }
    })
    .catch( (err) => {
      console.log(err) 
    });
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
