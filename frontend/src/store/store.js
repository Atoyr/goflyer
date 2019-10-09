import Vue from 'vue';
// Ajax通信ライブラリ
import axios from 'axios';
//
// // Json取得のベースURL
const URL_BASE = '/v1/dataframe/BTC_JPY/'
//
// Vue.js のインスタンス
module.exports = new Vue({
  data: {
    // Jsonデータ格納用
    data_frame: []
  },
  methods: {
    // Ajax通信でJsonを取得し、特定のプロパティに格納する
    // 取得したら GET_AJAX_COMPLETE で通知する
    get_ajax(url, name) {
      return axios.get(URL_BASE + url)
        .then((res) => {
          Vue.set(this, name, res.data);
          this.$emit('GET_AJAX_COMPLETE');
        });
    },
    // プロパティ名を指定してデータを取得
    get_data(name) {
      return this.$data[name];
    }
  }
});
