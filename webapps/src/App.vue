<template>
  <div id="contents">
      <GChart
        type="CandlestickChart"
        :data="chartData"
        :options="chartOptions"
      />
    </div>
</template>

<script>
import { GChart } from 'vue-google-charts';
import fetch from 'node-fetch';

export default {
  components: {
    GChart
  },
  methods : {
    getCandles(){
      fetch('http://localhost:8080/candles/1')
        .then(res => {
          return res.json();
        })
        .then(data => {
          this.chartData = [["","","","",""]]
          for (const c of data) {
            let temp = []
            temp.push(c.time)
            temp.push(c.low)
            temp.push(c.open)
            temp.push(c.close)
            temp.push(c.high)
            this.chartData.push(temp)
          }
        })
    }
  },
  mounted() {
    let self = this
    this.intervalId = setInterval(function () {
      self.getCandles()
    }, 1000)
  },
  beforeDestroy() {
    clearInterval(this.intervalId)
  },
  data() {
    return {
      chartData: [ ],
      chartOptions: {
        title: 'Company Performance',
        subtitle: 'Sales'
      },
      intervalId : undefined
    };
  }
};
</script>
