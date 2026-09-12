<!-- Card3D.vue -->
<template>
  <div ref="container" class="card-3d-container">
    <div class="info-tip">✨ 3D 卡牌 | 金属质感反光 | 鼠标拖动旋转</div>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'
import * as THREE from 'three'
import { OrbitControls } from 'three/addons/controls/OrbitControls.js'

const container = ref(null)

// 保存需要清理的对象
let scene, camera, renderer, controls, animationId
let cardGroup, gem, ring, particles
let clock

// 程序化生成渐变环境贴图（用于金属反光）
function createGradientEnvMap() {
  const size = 64
  const canvas = document.createElement('canvas')
  canvas.width = size
  canvas.height = size
  const ctx = canvas.getContext('2d')

  const gradient = ctx.createLinearGradient(0, 0, size, size)
  gradient.addColorStop(0, '#aaccff')
  gradient.addColorStop(0.5, '#6688cc')
  gradient.addColorStop(1, '#223355')
  ctx.fillStyle = gradient
  ctx.fillRect(0, 0, size, size)

  // 高光点，增加反光变化
  ctx.fillStyle = '#ffffff'
  ctx.beginPath()
  ctx.arc(size * 0.8, size * 0.2, size * 0.1, 0, Math.PI * 2)
  ctx.fill()

  ctx.fillStyle = '#ffccaa'
  ctx.beginPath()
  ctx.arc(size * 0.2, size * 0.7, size * 0.08, 0, Math.PI * 2)
  ctx.fill()

  return canvas
}

// 构建圆角矩形几何体
function createRoundedRectGeometry(width, height, radius, depth, bevelSize = 0.03, bevelThickness = 0.03, curveSegments = 12) {
  const shape = new THREE.Shape()
  const x = -width / 2
  const y = -height / 2

  shape.moveTo(x + radius, y)
  shape.lineTo(x + width - radius, y)
  shape.quadraticCurveTo(x + width, y, x + width, y + radius)
  shape.lineTo(x + width, y + height - radius)
  shape.quadraticCurveTo(x + width, y + height, x + width - radius, y + height)
  shape.lineTo(x + radius, y + height)
  shape.quadraticCurveTo(x, y + height, x, y + height - radius)
  shape.lineTo(x, y + radius)
  shape.quadraticCurveTo(x, y, x + radius, y)

  const extrudeSettings = {
    depth,
    bevelEnabled: true,
    bevelSegments: 4,
    bevelSize,
    bevelThickness,
    curveSegments
  }

  const geometry = new THREE.ExtrudeGeometry(shape, extrudeSettings)
  geometry.center()
  return geometry
}

// 初始化场景
function initScene() {
  const el = container.value
  const width = el.clientWidth
  const height = el.clientHeight

  // 场景
  scene = new THREE.Scene()
  scene.background = new THREE.Color(0x0a0a1a)

  // 相机
  camera = new THREE.PerspectiveCamera(45, width / height, 0.1, 1000)
  camera.position.set(2, 1.5, 4)
  camera.lookAt(0, 0, 0)

  // 渲染器
  renderer = new THREE.WebGLRenderer({
    antialias: true,
    alpha: false,
    powerPreference: 'high-performance'
  })
  renderer.setSize(width, height)
  renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2))
  renderer.shadowMap.enabled = true
  renderer.shadowMap.type = THREE.PCFSoftShadowMap
  renderer.outputEncoding = THREE.sRGBEncoding
  renderer.toneMapping = THREE.ACESFilmicToneMapping
  renderer.toneMappingExposure = 1.2
  el.appendChild(renderer.domElement)

  // 控制器
  controls = new OrbitControls(camera, renderer.domElement)
  controls.enableDamping = true
  controls.dampingFactor = 0.05
  controls.autoRotate = true
  controls.autoRotateSpeed = 1.8
  controls.enableZoom = true
  controls.zoomSpeed = 0.8
  controls.enablePan = false
  controls.target.set(0, 0, 0)
  controls.maxPolarAngle = Math.PI / 2.2
  controls.minDistance = 2.5
  controls.maxDistance = 8

  // 环境贴图（金属反光的关键）
  const envCanvas = createGradientEnvMap()
  const envTexture = new THREE.CanvasTexture(envCanvas)
  envTexture.mapping = THREE.EquirectangularReflectionMapping
  scene.environment = envTexture

  // 灯光
  const ambientLight = new THREE.AmbientLight(0x404060)
  scene.add(ambientLight)

  const mainLight = new THREE.PointLight(0xffeedd, 1.2, 20)
  mainLight.position.set(2, 3, 4)
  mainLight.castShadow = true
  mainLight.shadow.mapSize.width = 1024
  mainLight.shadow.mapSize.height = 1024
  scene.add(mainLight)

  const leftLight = new THREE.PointLight(0xaaccff, 0.8)
  leftLight.position.set(-3, 1, 2)
  scene.add(leftLight)

  const backLight = new THREE.PointLight(0xffaa88, 0.6)
  backLight.position.set(1, 0.5, -4)
  scene.add(backLight)

  const bottomLight = new THREE.PointLight(0x99aaff, 0.4)
  bottomLight.position.set(0, -2, 1)
  scene.add(bottomLight)

  const fillLight = new THREE.PointLight(0xffffff, 0.3)
  fillLight.position.set(1, 2, 3)
  scene.add(fillLight)

  // 构建卡牌组
  cardGroup = new THREE.Group()

  const widthCard = 1.8
  const heightCard = 2.6
  const depthCard = 0.08
  const radiusCard = 0.15

  // 卡牌主体
  const geometry = createRoundedRectGeometry(widthCard, heightCard, radiusCard, depthCard)
  const material = new THREE.MeshStandardMaterial({
    color: 0xd0d8e8,
    metalness: 0.92,
    roughness: 0.25,
    envMapIntensity: 1.2,
    side: THREE.DoubleSide
  })
  const cardMesh = new THREE.Mesh(geometry, material)
  cardMesh.castShadow = true
  cardMesh.receiveShadow = true
  cardGroup.add(cardMesh)

  // 正面边框
  const borderShape = new THREE.Shape()
  const bw = widthCard * 0.88
  const bh = heightCard * 0.88
  const br = radiusCard * 0.8

  borderShape.moveTo(-bw / 2 + br, -bh / 2)
  borderShape.lineTo(bw / 2 - br, -bh / 2)
  borderShape.quadraticCurveTo(bw / 2, -bh / 2, bw / 2, -bh / 2 + br)
  borderShape.lineTo(bw / 2, bh / 2 - br)
  borderShape.quadraticCurveTo(bw / 2, bh / 2, bw / 2 - br, bh / 2)
  borderShape.lineTo(-bw / 2 + br, bh / 2)
  borderShape.quadraticCurveTo(-bw / 2, bh / 2, -bw / 2, bh / 2 - br)
  borderShape.lineTo(-bw / 2, -bh / 2 + br)
  borderShape.quadraticCurveTo(-bw / 2, -bh / 2, -bw / 2 + br, -bh / 2)

  const borderGeo = new THREE.ExtrudeGeometry(borderShape, {
    depth: 0.02,
    bevelEnabled: true,
    bevelSegments: 2,
    bevelSize: 0.01,
    bevelThickness: 0.01,
    curveSegments: 8
  })
  borderGeo.center()

  const borderMat = new THREE.MeshStandardMaterial({
    color: 0xffdd99,
    metalness: 0.95,
    roughness: 0.2,
    emissive: 0x332200,
    emissiveIntensity: 0.2,
    envMapIntensity: 1.5
  })
  const borderMesh = new THREE.Mesh(borderGeo, borderMat)
  borderMesh.position.z = depthCard / 2 + 0.01
  borderMesh.castShadow = true
  cardGroup.add(borderMesh)

  // 中心宝石
  const gemGeo = new THREE.OctahedronGeometry(0.25, 0)
  const gemMat = new THREE.MeshStandardMaterial({
    color: 0x66ccff,
    metalness: 1.0,
    roughness: 0.15,
    emissive: 0x224466,
    emissiveIntensity: 0.8,
    envMapIntensity: 1.8
  })
  gem = new THREE.Mesh(gemGeo, gemMat)
  gem.position.z = depthCard / 2 + 0.08
  gem.castShadow = true
  cardGroup.add(gem)

  // 宝石光环
  const ringGeo = new THREE.TorusGeometry(0.35, 0.02, 16, 32)
  const ringMat = new THREE.MeshStandardMaterial({
    color: 0xaaddff,
    metalness: 0.9,
    roughness: 0.3,
    emissive: 0x3366aa,
    emissiveIntensity: 0.5
  })
  ring = new THREE.Mesh(ringGeo, ringMat)
  ring.position.z = depthCard / 2 + 0.06
  ring.rotation.x = Math.PI / 2
  ring.rotation.z = 0.3
  cardGroup.add(ring)

  // 背面装饰板
  const backPlateGeo = new THREE.BoxGeometry(widthCard * 0.92, heightCard * 0.92, 0.01)
  const backPlateMat = new THREE.MeshStandardMaterial({
    color: 0x334466,
    metalness: 0.8,
    roughness: 0.4,
    emissive: 0x112233,
    emissiveIntensity: 0.1
  })
  const backPlate = new THREE.Mesh(backPlateGeo, backPlateMat)
  backPlate.position.z = -depthCard / 2 - 0.01
  backPlate.castShadow = true
  cardGroup.add(backPlate)

  scene.add(cardGroup)

  // 半透明接收阴影的地面
  const groundGeo = new THREE.CircleGeometry(5, 32)
  const groundMat = new THREE.MeshStandardMaterial({
    color: 0x1a1a2e,
    transparent: true,
    opacity: 0.15,
    side: THREE.DoubleSide,
    metalness: 0.5,
    roughness: 0.8
  })
  const ground = new THREE.Mesh(groundGeo, groundMat)
  ground.rotation.x = -Math.PI / 2
  ground.position.y = -heightCard / 2 - 0.5
  ground.receiveShadow = true
  scene.add(ground)

  // 粒子氛围
  const particleCount = 300
  const particlesGeo = new THREE.BufferGeometry()
  const particlePositions = new Float32Array(particleCount * 3)
  for (let i = 0; i < particleCount * 3; i += 3) {
    particlePositions[i] = (Math.random() - 0.5) * 20
    particlePositions[i + 1] = (Math.random() - 0.5) * 20
    particlePositions[i + 2] = (Math.random() - 0.5) * 20
  }
  particlesGeo.setAttribute('position', new THREE.BufferAttribute(particlePositions, 3))
  const particlesMat = new THREE.PointsMaterial({
    color: 0xaaccff,
    size: 0.05,
    transparent: true,
    opacity: 0.5,
    blending: THREE.AdditiveBlending,
    depthWrite: false
  })
  particles = new THREE.Points(particlesGeo, particlesMat)
  scene.add(particles)

  clock = new THREE.Clock()

  // 监听窗口尺寸变化
  window.addEventListener('resize', onWindowResize)
}

// 动画循环
function animate() {
  animationId = requestAnimationFrame(animate)

  const elapsedTime = performance.now() * 0.001

  controls.update()

  if (gem) {
    gem.rotation.y = elapsedTime * 0.5
    gem.rotation.x = Math.sin(elapsedTime * 0.8) * 0.2
  }
  if (ring) {
    ring.rotation.z += 0.005
    ring.rotation.x = Math.PI / 2 + Math.sin(elapsedTime * 0.6) * 0.1
  }
  if (particles) {
    particles.rotation.y += 0.0001
  }

  renderer.render(scene, camera)
}

// 窗口自适应
function onWindowResize() {
  if (!container.value || !camera || !renderer) return
  const el = container.value
  const width = el.clientWidth
  const height = el.clientHeight

  camera.aspect = width / height
  camera.updateProjectionMatrix()
  renderer.setSize(width, height)
  renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2))
}

// 清理资源
function disposeScene() {
  if (animationId) cancelAnimationFrame(animationId)
  window.removeEventListener('resize', onWindowResize)

  if (controls) controls.dispose()

  if (scene) {
    scene.traverse((obj) => {
      if (obj.geometry) obj.geometry.dispose()
      if (obj.material) {
        if (Array.isArray(obj.material)) {
          obj.material.forEach((m) => m.dispose())
        } else {
          obj.material.dispose()
        }
      }
    })
  }

  if (renderer) {
    renderer.dispose()
    renderer.forceContextLoss()
    if (renderer.domElement && renderer.domElement.parentNode) {
      renderer.domElement.parentNode.removeChild(renderer.domElement)
    }
  }

  scene = null
  camera = null
  renderer = null
  controls = null
  cardGroup = null
  gem = null
  ring = null
  particles = null
}

onMounted(() => {
  initScene()
  animate()
})

onBeforeUnmount(() => {
  disposeScene()
})
</script>

<style scoped>
.card-3d-container {
  position: relative;
  width: 100%;
  height: 100%;
  min-height: 500px;
  overflow: hidden;
  background-color: #0a0a1a;
}

.info-tip {
  position: absolute;
  bottom: 20px;
  left: 20px;
  color: rgba(255, 255, 255, 0.6);
  font-size: 14px;
  pointer-events: none;
  text-shadow: 0 0 10px rgba(0, 0, 0, 0.8);
  z-index: 100;
}
</style>