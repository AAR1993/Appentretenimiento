// Variables globales para controlar la búsqueda
let timeoutId;
const delay = 300; // 300ms de delay antes de hacer la búsqueda

// Función para cargar todos los lugares
async function cargarLugares() {
    try {
        const response = await fetch('https://appentretenimiento-1.onrender.com/api/lugares');
        if (!response.ok) {
            throw new Error('Error al obtener los lugares');
        }

        const lugares = await response.json();
        mostrarLugares(lugares);
    } catch (error) {
        console.error(error);
        mostrarError('Error al cargar los lugares');
    }
}

// Función para mostrar los lugares en la lista
function mostrarLugares(lugares) {
    const lista = document.getElementById('lista-lugares');
    lista.innerHTML = ''; // Limpiar lista existente

    if (lugares.length === 0) {
        lista.innerHTML = '<p>No se encontraron lugares que coincidan con la búsqueda.</p>';
        return;
    }

    lugares.forEach(lugar => {
        const div = document.createElement('div');
        div.className = 'lugar-item';
        div.innerHTML = `
            <h3>${lugar.Nombre || 'Nombre no disponible'}</h3>
            <img src="${lugar.Imagen || 'ruta/a/la/imagen/predeterminada.jpg'}" alt="${lugar.Nombre || 'Imagen no disponible'}" class="lugar-image" />
            <p>${lugar.Descripcion || 'Descripción no disponible'}</p>
            <p><strong>Horario:</strong> ${lugar.Horario || 'Horario no disponible'}</p>
        `;
        lista.appendChild(div);
    });
}

// Función para mostrar mensajes de error
function mostrarError(mensaje) {
    const lista = document.getElementById('lista-lugares');
    lista.innerHTML = `<p class="error">${mensaje}</p>`;
}

// Función principal de búsqueda
function buscarLugares(query) {
    // Limpiar timeout anterior si existe
    if (timeoutId) {
        clearTimeout(timeoutId);
    }

    // Establecer nuevo timeout
    timeoutId = setTimeout(async () => {
        try {
            const response = await fetch('https://appentretenimiento-1.onrender.com/api/lugares');
            if (!response.ok) {
                throw new Error('Error al obtener los lugares');
            }

            const lugares = await response.json();
            // Convertir la query a un array de palabras
            const palabras = query.toLowerCase().split(' ').filter(Boolean);
            
            // Función para verificar si un texto contiene todas las palabras
            const contieneTodasLasPalabras = (texto) => {
                const textoLower = texto.toLowerCase();
                return palabras.every(palabra => 
                    textoLower.includes(palabra) ||
                    // También buscar parcialmente al inicio de palabras
                    textoLower.split(' ').some(palabraTexto => 
                        palabraTexto.startsWith(palabra)
                    )
                );
            };

            // Filtrar lugares que contengan todas las palabras en nombre o descripción
            const resultados = lugares.filter(lugar => 
                contieneTodasLasPalabras(lugar.Nombre) ||
                contieneTodasLasPalabras(lugar.Descripcion)
            );
            mostrarLugares(resultados);
        } catch (error) {
            console.error(error);
            mostrarError('Error al realizar la búsqueda');
        }
    }, delay);
}

// Inicialización de la aplicación
document.addEventListener('DOMContentLoaded', () => {
    // Cargar lugares al iniciar
    cargarLugares();
    
    // Configurar búsqueda automática
    const searchInput = document.getElementById('search-input');
    if (searchInput) {
        searchInput.addEventListener('input', (e) => {
            const query = e.target.value;
            buscarLugares(query);
        });
    }
});
