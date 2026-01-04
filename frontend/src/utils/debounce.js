/**
 * 防抖函数 (Debounce)
 * 
 * 作用：限制函数在短时间内的高频触发。只有当 wait 毫秒内没有新的调用时，
 *       才会真正执行一次 func。
 * 
 * 场景：用户连续输入时（文本框 input 事件），只在停止输入后执行搜索或上报，
 *       避免过于频繁的 API 请求或计算。
 * 
 * @param {Function} func - 需要防抖执行的原始函数
 * @param {Number} wait   - 等待时间 (毫秒)
 * @returns {Function}    - 返回一个新的封装好的防抖函数
 */
export function debounce(func, wait) {
    let timeout;
    return function (...args) {
        const context = this;
        // 如果之前有定时器（说明还在等待期间又触发了），则清除它，重新计时
        clearTimeout(timeout);
        // 设置新的定时器
        timeout = setTimeout(() => func.apply(context, args), wait);
    };
}
