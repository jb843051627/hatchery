const api = '/api';
const content = document.getElementById('content');
document.querySelectorAll('.tab').forEach(btn => {
  btn.addEventListener('click', () => {
    document.querySelectorAll('.tab').forEach(b => b.classList.remove('active'));
    btn.classList.add('active');
    loadView(btn.dataset.view);
  });
});

async function fetchJSON(url) {
  const res = await fetch(url);
  if (!res.ok) throw new Error(`${res.status} ${res.statusText}`);
  return res.json();
}

async function loadView(view) {
  try {
    if (view === 'dashboard') return await loadDashboard();
    if (view === 'incubators') return await loadIncubators();
    if (view === 'batches') return await loadBatches();
    if (view === 'readings') return await loadReadings();
    if (view === 'alerts') return await loadAlerts();
  } catch (e) {
    content.innerHTML = `<div class="card"><p>加载失败: ${e.message}</p></div>`;
  }
}

async function loadDashboard() {
  const d = await fetchJSON(`${api}/reports/dashboard`);
  content.innerHTML = `
    <div class="stat-grid">
      <div class="card stat"><div class="num">${d.total_incubators}</div><div class="label">孵化器总数</div></div>
      <div class="card stat"><div class="num">${d.active_batches}</div><div class="label">活跃批次</div></div>
      <div class="card stat"><div class="num">${d.pending_maintenance}</div><div class="label">待维护</div></div>
    </div>
    <div class="card"><h2>近期出雏记录</h2>
      <table><thead><tr><th>ID</th><th>批次</th><th>出雏数</th><th>健康</th><th>弱雏</th><th>死胚</th></tr></thead>
      <tbody>${(d.recent_hatch_records||[]).map(r=>`<tr><td>${r.id}</td><td>${r.batch_id}</td><td>${r.hatched_count}</td><td>${r.healthy_count}</td><td>${r.weak_count}</td><td>${r.dead_count}</td></tr>`).join('')}</tbody>
    </table></div>`;
}

async function loadIncubators() {
  const list = await fetchJSON(`${api}/incubators`);
  content.innerHTML = `<div class="card"><h2>孵化器列表</h2>
    <table><thead><tr><th>ID</th><th>名称</th><th>容量</th><th>状态</th><th>位置</th></tr></thead>
    <tbody>${(list||[]).map(i=>`<tr><td>${i.id}</td><td>${i.name}</td><td>${i.capacity}</td><td><span class="badge ${i.status}">${i.status}</span></td><td>${i.location}</td></tr>`).join('')}</tbody>
  </table></div>`;
}

async function loadBatches() {
  const list = await fetchJSON(`${api}/batches?status=incubating`);
  content.innerHTML = `<div class="card"><h2>孵化中批次</h2>
    <table><thead><tr><th>ID</th><th>孵化器</th><th>蛋数</th><th>品种</th><th>状态</th><th>预计出雏</th></tr></thead>
    <tbody>${(list||[]).map(b=>`<tr><td>${b.id}</td><td>${b.incubator_id}</td><td>${b.egg_count}</td><td>${b.species}</td><td><span class="badge incubating">${b.status}</span></td><td>${b.expected_hatch_date.slice(0,10)}</td></tr>`).join('')}</tbody>
  </table></div>`;
}

async function loadReadings() {
  content.innerHTML = `<div class="card"><h2>传感器读数</h2><p>请在概览页面查看各孵化器最新读数。</p></div>`;
}

async function loadAlerts() {
  content.innerHTML = `<div class="card"><h2>告警列表</h2><p>请在各孵化器详情页面查看告警。</p></div>`;
}

loadView('dashboard');
